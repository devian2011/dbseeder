package seeder

import (
	"dbseeder/internal/seeder/providers"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"dbseeder/internal/modifiers"
	"dbseeder/internal/schema"
	"dbseeder/internal/seeder/generators/dependence"
	"dbseeder/internal/seeder/generators/fake"
	"dbseeder/internal/seeder/generators/list"
	"dbseeder/pkg/helper"
)

type DbProvider interface {
	GetAll(tableName string) ([]map[string]any, error)
	Truncate(tableName string) error
	Insert(tableName string, columns []string, values []any) error
	Commit() error
	Rollback() error
}

type RelationValues struct {
	//TODO: Optimize store values
	Tables map[string][]map[string]any
}

func (t *RelationValues) Add(tableCode string, rows []map[string]any) {
	t.Tables[tableCode] = rows
}

func (t *RelationValues) Get(tableCode string) ([]map[string]any, error) {
	if v, exists := t.Tables[tableCode]; exists {
		return v, nil
	}
	return nil, fmt.Errorf("values for table: %s didn't be generated", tableCode)
}

func (t *RelationValues) isTableGenerated(tableCode string) bool {
	if _, exists := t.Tables[tableCode]; exists {
		return true
	}
	return false
}

type connectionPool struct {
	connections map[string]DbProvider
}

func newConnectionPool() *connectionPool {
	return &connectionPool{
		connections: make(map[string]DbProvider, 1),
	}
}

func (pool *connectionPool) initConnection(code, driver, dsn string) error {
	var err error
	if _, exists := pool.connections[code]; exists {
		return fmt.Errorf("database with code name: %s already exists", code)
	}
	switch driver {
	case "pgx":
		pool.connections[code], err = providers.NewPsqlProvider(dsn)
		return err
	case "mysql":
		pool.connections[code], err = providers.NewMysqlProvider(dsn)
		return err
	default:
		return fmt.Errorf("unsupported database driver: %s for code: %s", driver, code)
	}
}

func (pool *connectionPool) getConnection(db string) (DbProvider, error) {
	if trx, exists := pool.connections[db]; exists {
		return trx, nil
	}

	return nil, fmt.Errorf("unknown transaction for db: %s", db)
}

func (pool *connectionPool) commitTransactions() {
	for code, conn := range pool.connections {
		conn.Commit()
		logrus.Infof("Transaction for %s has been commited", code)
	}
}

func (pool *connectionPool) rollbackTransactions() {
	for code, conn := range pool.connections {
		conn.Rollback()
		logrus.Infof("Transaction for %s has been rollbacked", code)
	}
}

type Seeder struct {
	connPool       *connectionPool
	schema         *schema.Schema
	relValues      *RelationValues
	tableGenerator *TableGenerator
}

func NewSeeder(sch *schema.Schema, plugins *modifiers.ModifierStore) *Seeder {
	relValues := &RelationValues{
		Tables: make(map[string][]map[string]any, 0),
	}
	return &Seeder{
		connPool:       newConnectionPool(),
		schema:         sch,
		tableGenerator: &TableGenerator{relationValues: relValues, plugins: plugins},
		relValues:      relValues,
	}
}

func (seeder *Seeder) initConnections() error {
	for _, d := range seeder.schema.Databases.Databases {
		err := seeder.connPool.initConnection(d.Name, d.Driver, d.DSN)
		if err != nil {
			return err
		}
	}

	return nil
}

func (seeder *Seeder) Run() error {
	initConnectionErr := seeder.initConnections()
	if initConnectionErr != nil {
		return initConnectionErr
	}
	// Rollback transactions and close connection
	defer seeder.connPool.rollbackTransactions()

	err := seeder.schema.Tree.Walk(seeder.walkingFn)

	if err != nil {
		seeder.connPool.rollbackTransactions()
	} else {
		seeder.connPool.commitTransactions()
	}

	return err
}

func (seeder *Seeder) walkingFn(code string, node *schema.TableDependenceNode) error {
	if seeder.relValues.isTableGenerated(code) {
		return nil
	}
	/// If table has dependencies - truncate it
	if len(node.Dependencies) > 0 && node.Table.Action == schema.ActionGenerate {
		truncateErr := seeder.truncate(node.DbName, node.Table.Name)
		if truncateErr != nil {
			return truncateErr
		}
	}

	for _, dependency := range node.Dependencies {
		err := seeder.walkingFn(dependency.Code, dependency)
		if err != nil {
			return err
		}
	}

	switch node.Table.Action {
	case schema.ActionGenerate:
		td, tableGenErr := seeder.tableGenerator.generate(code, node)
		if tableGenErr != nil {
			return tableGenErr
		}

		insertErr := seeder.insert(node.DbName, node.Table.Name, td.columns, td.values)
		if insertErr != nil {
			return insertErr
		}
		/// Mysql and Postgres have different ways for get ids, that's why simple way is select all generated data
		return seeder.loadDataFromDb(code, node)
	case schema.ActionGet:
		return seeder.loadDataFromDb(code, node)
	default:
		return fmt.Errorf("unknown action for table: %s", code)
	}
}

type TableGenerator struct {
	relationValues *RelationValues
	plugins        *modifiers.ModifierStore
}

type tableData struct {
	rowsCount      int
	columns        []string
	values         []any
	rowsHashes     map[string]bool
	relations      map[string]map[int]bool
	relationValues *RelationValues
	node           *schema.TableDependenceNode
	plugins        *modifiers.ModifierStore
}

type orderedColumns struct {
	cols []string
}

func (o *orderedColumns) exists(val string) bool {
	for _, existValues := range o.cols {
		if existValues == val {
			return true
		}
	}

	return false
}

// / TODO: Check circular dependencies
func getOrderedColumns(columns *orderedColumns, fieldName string, fields map[string]schema.Field) {
	field := fields[fieldName]
	if field.IsExpressionDependence() {
		for _, r := range field.Depends.Expression.Rows {
			if columns.exists(r) {
				continue
			}
			getOrderedColumns(columns, r, fields)
		}
	}
	columns.cols = append(columns.cols, fieldName)
}

func (generator *TableGenerator) getColumns(node *schema.TableDependenceNode) ([]string, error) {
	columns := &orderedColumns{cols: make([]string, 0)}
	for fieldName, fldVal := range node.Table.Fields {
		if fldVal.Generation == schema.GenerationTypeDb {
			continue
		}
		if columns.exists(fieldName) {
			continue
		}

		getOrderedColumns(columns, fieldName, node.Table.Fields)
	}

	return columns.cols, nil
}

func (generator *TableGenerator) initTableData(code string, node *schema.TableDependenceNode, relValues *RelationValues) (*tableData, error) {
	/// Found length for generated values
	rowsCount := node.Table.GetRowsCount()
	if rowsCount == 0 {
		return nil, fmt.Errorf(
			"count rows and fill part for table: %s is empty, set count generated rows or make fill data",
			code)
	}
	columns, err := generator.getColumns(node)
	if err != nil {
		return nil, err
	}

	return &tableData{
		rowsCount:      rowsCount,
		columns:        columns,
		values:         make([]any, 0, rowsCount*len(columns)),
		rowsHashes:     make(map[string]bool, rowsCount),
		relations:      make(map[string]map[int]bool, 0),
		relationValues: relValues,
		node:           node,
		plugins:        generator.plugins,
	}, nil
}

func (generator *TableGenerator) generate(code string, node *schema.TableDependenceNode) (*tableData, error) {
	td, tableDataErr := generator.initTableData(code, node, generator.relationValues)
	if tableDataErr != nil {
		return nil, tableDataErr
	}

	for c := 0; c < td.rowsCount; c++ {
		exists := true
		for ok := true; ok; ok = exists {
			rowValues, rowGeneratedErr := td.generateRow(td.columns, c)
			if rowGeneratedErr != nil {
				return nil, rowGeneratedErr
			}

			// It table require no duplicates - check that this values has no identical rows
			if td.node.Table.NoDuplicates {
				sliceHash := helper.SliceHash(rowValues)
				if _, exists = td.rowsHashes[sliceHash]; !exists {
					td.values = append(td.values, rowValues...)
					td.rowsHashes[sliceHash] = true
				}
			} else {
				td.values = append(td.values, rowValues...)
				exists = false
			}
		}
	}

	return td, nil
}

func (generator *tableData) generateRow(columns []string, rowNum int) ([]any, error) {
	rowValues := make(map[string]any, len(columns))
	values := make([]any, 0, len(columns))
	for _, fieldName := range columns {
		if len(generator.node.Table.Fill) > rowNum {
			if v, exists := generator.node.Table.Fill[rowNum][fieldName]; exists {
				values = append(values, v)
				rowValues[fieldName] = v
				continue
			}
		}

		v, err := generator.generateFieldData(fieldName, rowValues)

		if err == ErrGetFromDb {
			continue
		}

		if err != nil {
			return nil, err
		}

		values = append(values, v)
		rowValues[fieldName] = v

	}

	return values, nil
}

var ErrGetFromDb = errors.New("get field val from db")

func (generator *tableData) generateFieldData(fieldName string, rowValues map[string]any) (any, error) {
	fieldValue := generator.node.Table.Fields[fieldName]
	switch fieldValue.Generation {
	case schema.GenerationTypeFaker:
		vl, err := fake.Generate(fieldName, fieldValue)
		if err != nil {
			return vl, err
		}
		return generator.plugins.ApplyList(fieldValue.Plugins, vl)
	case schema.GenerationTypeList:
		vl, err := list.Generate(fieldName, fieldValue)
		if err != nil {
			return vl, err
		}
		return generator.plugins.ApplyList(fieldValue.Plugins, vl)
	case schema.GenerationDepends:
		if fieldValue.IsFkDependence() {
			return dependence.GenerateForeign(fieldValue, generator.relationValues, generator.relations)
		}
		if fieldValue.IsExpressionDependence() {
			return dependence.GenerateExpression(fieldValue, rowValues)
		}
		return nil, fmt.Errorf("unknown dependence field generation type for %s", fieldName)
	case schema.GenerationTypeDb:
		return nil, ErrGetFromDb
	default:
		return nil, fmt.Errorf("unknown field generation type for %s", fieldName)
	}
}

// Db functions
func (seeder *Seeder) loadDataFromDb(code string, node *schema.TableDependenceNode) error {
	if node.HasDependents() {
		conn, err := seeder.connPool.getConnection(node.DbName)
		if err != nil {
			return err
		}

		result, err := conn.GetAll(node.Table.Name)
		if err != nil {
			return nil
		}
		seeder.relValues.Add(code, result)
	}
	return nil
}

func (seeder *Seeder) insert(db string, tableName string, columns []string, values []any) error {
	conn, err := seeder.connPool.getConnection(db)
	if err != nil {
		return err
	}
	return conn.Insert(tableName, columns, values)
}

func (seeder *Seeder) truncate(db, tableName string) error {
	conn, err := seeder.connPool.getConnection(db)
	if err != nil {
		return err
	}
	return conn.Truncate(tableName)
}
