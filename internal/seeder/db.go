package seeder

import (
	"dbseeder/internal/schema"
	"dbseeder/internal/seeder/providers"
	"fmt"
	"github.com/sirupsen/logrus"
)

type DbProvider interface {
	GetAll(tableName string) ([]map[string]any, error)
	Truncate(tableName string) error
	Insert(tableName string, columns []string, values []any) error
	Commit() error
	Rollback() error
}

type dbConnPool struct {
	connections map[string]DbProvider
}

func newConnectionPool() *dbConnPool {
	return &dbConnPool{
		connections: make(map[string]DbProvider, 1),
	}
}

func (pool *dbConnPool) initConnection(code, driver, dsn string) error {
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

func (pool *dbConnPool) getConnection(db string) (DbProvider, error) {
	if trx, exists := pool.connections[db]; exists {
		return trx, nil
	}

	return nil, fmt.Errorf("unknown transaction for db: %s", db)
}

func (pool *dbConnPool) commit() {
	for code, conn := range pool.connections {
		conn.Commit()
		logrus.Infof("Transaction for %s has been commited", code)
	}
}

func (pool *dbConnPool) rollback() {
	for code, conn := range pool.connections {
		conn.Rollback()
		logrus.Infof("Transaction for %s has been rollbacked", code)
	}
}

type Db struct {
	pool *dbConnPool
}

func NewDb(databases map[string]*schema.Database) (*Db, error) {
	pool := &dbConnPool{connections: make(map[string]DbProvider, len(databases))}
	for _, d := range databases {
		err := pool.initConnection(d.Name, d.Driver, d.DSN)
		if err != nil {
			pool.rollback()
			return nil, err
		}
	}

	return &Db{pool: pool}, nil
}

func (db *Db) loadDataFromDb(code, tableName string) ([]map[string]any, error) {
	conn, err := db.pool.getConnection(code)
	if err != nil {
		return nil, err
	}

	return conn.GetAll(tableName)
}

func (db *Db) insert(dbCode string, tableName string, columns []string, values []any) error {
	conn, err := db.pool.getConnection(dbCode)
	if err != nil {
		return err
	}
	return conn.Insert(tableName, columns, values)
}

func (db *Db) truncate(dbCode, tableName string) error {
	conn, err := db.pool.getConnection(dbCode)
	if err != nil {
		return err
	}
	return conn.Truncate(tableName)
}
