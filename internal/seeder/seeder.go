package seeder

import (
	"dbseeder/internal/schema"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Seeder struct {
	connections     map[string]*sqlx.DB
	trxs            map[string]*sqlx.Tx
	schema          *schema.Schema
	generatedTables map[string]bool
	generationErrCh error
}

type TableValues struct {
	//TODO: Optimize store values
	Tables map[string][]map[string]any
}

func (t *TableValues) Add(tableCode string, rows []map[string]any) {
	t.Tables[tableCode] = rows
}

func NewSeeder(sch *schema.Schema) (*Seeder, error) {
	seeder := &Seeder{
		connections:     make(map[string]*sqlx.DB),
		trxs:            make(map[string]*sqlx.Tx),
		schema:          sch,
		generatedTables: make(map[string]bool),
		generationErrCh: nil,
	}

	return seeder, seeder.initConnections()
}

func (seeder *Seeder) initConnections() error {
	var err error
	for _, d := range seeder.schema.Databases.Databases {
		if _, exists := seeder.connections[d.Name]; exists {
			return errors.New(fmt.Sprintf("database with code name: %s already exists", d.Name))
		}
		seeder.connections[d.Name], err = sqlx.Open("pgx", d.DSN)
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
	seeder.schema.Tree.Walk(seeder.walkingFn)

	if seeder.generationErrCh != nil {
		seeder.rollbackTrxs()
	} else {
		seeder.commitTrxs()
	}

	return seeder.generationErrCh
}

func (seeder *Seeder) walkingFn(code string, node *schema.TableDependenceNode) {
	if generated, exists := seeder.generatedTables[code]; generated && exists {
		return
	}
	for _, dependency := range node.Dependencies {
		if generated, exists := seeder.generatedTables[dependency.Code]; !generated || !exists {
			seeder.walkingFn(dependency.Code, dependency)
		}
	}
	
}
