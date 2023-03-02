package seeder

import (
	"dbseeder/internal/modifiers"
	"dbseeder/internal/schema"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type Seeder struct {
	connections     map[string]*sqlx.DB
	trxs            map[string]*sqlx.Tx
	schema          *schema.Schema
	plugins         *modifiers.ModifierStore
	generatedValues *TableValues
}

type TableValues struct {
	//TODO: Optimize store values
	Tables map[string][]map[string]any
}

func (t *TableValues) Add(tableCode string, rows []map[string]any) {
	t.Tables[tableCode] = rows
}

func (t *TableValues) Get(tableCode string) ([]map[string]any, error) {
	if v, exists := t.Tables[tableCode]; exists {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("values for table: %s didn't be generated", tableCode))
}

func (t *TableValues) IsTableGenerated(tableCode string) bool {
	if _, exists := t.Tables[tableCode]; exists {
		return true
	}
	return false
}

func NewSeeder(sch *schema.Schema, plugins *modifiers.ModifierStore) (*Seeder, error) {
	seeder := &Seeder{
		connections: make(map[string]*sqlx.DB),
		trxs:        make(map[string]*sqlx.Tx),
		schema:      sch,
		generatedValues: &TableValues{
			Tables: make(map[string][]map[string]any, 0),
		},
		plugins: plugins,
	}

	return seeder, seeder.initConnections()
}

func (seeder *Seeder) initConnections() error {
	var err error
	for _, d := range seeder.schema.Databases.Databases {
		if _, exists := seeder.connections[d.Name]; exists {
			return errors.New(fmt.Sprintf("database with code name: %s already exists", d.Name))
		}
		seeder.connections[d.Name], err = sqlx.Open(d.Driver, d.DSN)
		if err != nil {
			return errors.New(fmt.Sprintf("open connection for %s has error: %s", d.Name, err.Error()))
		}
		if err = seeder.connections[d.Name].Ping(); err != nil {
			return errors.New(fmt.Sprintf("error ping for %s error: %s", d.Name, err.Error()))
		}
		seeder.connections[d.Name].SetMaxOpenConns(10)
		seeder.connections[d.Name].SetMaxIdleConns(10)
	}

	return nil
}

func (seeder *Seeder) startTrxs() error {
	var initTxErr error
	for code, conn := range seeder.connections {
		seeder.trxs[code], initTxErr = conn.Beginx()
		if initTxErr != nil {
			logrus.Errorf("Error on start transaction for %s. Err: %s", code, initTxErr)
			seeder.rollbackTrxs()
			return initTxErr
		}
		logrus.Infof("Start transaction for %s", code)
	}

	return nil
}

func (seeder *Seeder) commitTrxs() {
	for code, trx := range seeder.trxs {
		trx.Commit()
		logrus.Infof("Transaction for %s has been commited", code)
	}
}

func (seeder *Seeder) rollbackTrxs() {
	for code, trx := range seeder.trxs {
		trx.Rollback()
		logrus.Infof("Transaction for %s has been rollbacked", code)
	}
}

func (seeder *Seeder) closeConnections() {
	for _, c := range seeder.connections {
		c.Close()
	}
}

func (seeder *Seeder) Run() error {
	defer seeder.closeConnections()
	if startTrxsError := seeder.startTrxs(); startTrxsError != nil {
		seeder.rollbackTrxs()
	}
	err := seeder.schema.Tree.Walk(seeder.walkingFn)

	if err != nil {
		seeder.rollbackTrxs()
	} else {
		seeder.commitTrxs()
	}

	return err
}

func (seeder *Seeder) walkingFn(code string, node *schema.TableDependenceNode) error {
	if seeder.generatedValues.IsTableGenerated(code) {
		return nil
	}
	for _, dependency := range node.Dependencies {
		err := seeder.walkingFn(dependency.Code, dependency)
		if err != nil {
			return err
		}
	}

	if node.Table.Action == schema.ActionGenerate {
		/// Found length for generated values
		rowsCount := node.Table.GetRowsCount()
		if rowsCount == 0 {
			return errors.New(
				fmt.Sprintf(
					"count rows and fill part for table: %s is empty, set count generated rows or make fill data",
					code))
		}

		// Drop data from db
		truncateErr := seeder.truncate(node.DbName, node.Table.Name)
		if truncateErr != nil {
			return truncateErr
		}

		// Columns for insert to db
		columns := make([]string, 0, len(node.Table.Fields))
		for fieldName, fVal := range node.Table.Fields {
			if fVal.Generation == schema.GenerationTypeDb {
				continue
			}
			columns = append(columns, fieldName)
		}
		// Values for insert to db
		values := make([]any, 0, rowsCount*len(columns))

		for c := 0; c < rowsCount; c++ {
			for _, fieldName := range columns {
				if len(node.Table.Fill) > c {
					if v, exists := node.Table.Fill[c][fieldName]; exists {
						values = append(values, v)
						continue
					}
				}

				fieldValue := node.Table.Fields[fieldName]

				switch fieldValue.Generation {
				case schema.GenerationTypeFaker:
					v, vErr := Generate(fieldValue.Type)
					if vErr != nil {
						return vErr
					}
					values = append(values, v)
				case schema.GenerationTypeList:
					rIndex := rand.Intn(len(fieldValue.List) - 1)
					values = append(values, fieldValue.List[rIndex])
				case schema.GenerationDepends:
					generatedVls, findErr := seeder.generatedValues.Get(
						fmt.Sprintf("%s.%s",
							fieldValue.Depends.Foreign[0].Db,
							fieldValue.Depends.Foreign[0].Table))
					if findErr != nil {
						return findErr
					}
					rIndex := rand.Intn(len(generatedVls) - 1)
					values = append(values, generatedVls[rIndex][fieldValue.Depends.Foreign[0].Field])
				}
			}
		}

		insertErr := seeder.insertToTable(node.DbName, node.Table.Name, columns, values)
		if insertErr != nil {
			return insertErr
		}

		/// Mysql and Postgres Has different ways for get ids, that's why simple way is select all generated data
		return seeder.getDataFromDb(code, node)
	}

	return seeder.getDataFromDb(code, node)
}

func (seeder *Seeder) getDataFromDb(code string, node *schema.TableDependenceNode) error {
	if node.HasDependents() {
		result := make([]map[string]any, 0)
		rows, err := seeder.trxs[node.DbName].Queryx(fmt.Sprintf("SELECT * FROM %s", node.Table.Name))
		defer rows.Close()
		if err != nil {
			return errors.New(
				fmt.Sprintf(
					"error for get values from table: %s. Error: %s",
					code,
					err.Error()))
		}
		for rows.Next() {
			val := make(map[string]any, 0)
			rows.MapScan(val)
			result = append(result, val)
		}

		seeder.generatedValues.Add(code, result)
	}
	return nil
}

func (seeder *Seeder) insertToTable(db string, tableName string, columns []string, values []any) error {
	if conn, exists := seeder.trxs[db]; exists {
		rowsCountForInsert := len(values) / len(columns)
		sql := conn.Rebind(generateInsertSql(tableName, columns, rowsCountForInsert))
		return conn.QueryRowx(sql, values...).Err()
	}

	return errors.New(fmt.Sprintf("unknown db with code: %s", db))
}

func (seeder *Seeder) truncate(db, tableName string) error {
	if conn, exists := seeder.trxs[db]; exists {
		return conn.QueryRowx(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Err()
	}

	return errors.New(fmt.Sprintf("unknown db with code: %s", db))
}
