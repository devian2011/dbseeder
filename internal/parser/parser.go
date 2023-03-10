package parser

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Databases map[string]Database `yaml:"databases" json:"databases"`
}

type Database struct {
	Driver string `json:"driver" yaml:"driver"`
	Name   string `json:"name" yaml:"name"` // Database name
	DSN    string `json:"dsn" yaml:"dsn"`   // Database DSN
}

type connectionPool struct {
	connections map[string]*sqlx.DB
}

func newConnectionPool() *connectionPool {
	return &connectionPool{
		connections: make(map[string]*sqlx.DB),
	}
}

func (pool *connectionPool) initConnection(code, driver, dsn string) error {
	var err error
	if _, exists := pool.connections[code]; exists {
		return errors.New(fmt.Sprintf("database with code name: %s already exists", code))
	}
	pool.connections[code], err = sqlx.Open(driver, dsn)
	if err != nil {
		return errors.New(fmt.Sprintf("open connection for %s has error: %s", code, err.Error()))
	}
	if err = pool.connections[code].Ping(); err != nil {
		return errors.New(fmt.Sprintf("error ping for %s error: %s", code, err.Error()))
	}
	pool.connections[code].SetMaxOpenConns(10)
	pool.connections[code].SetMaxIdleConns(10)
	return nil
}

type Parser struct {
	cfg      *Config
	connPool *connectionPool
}

func NewParser(cfg *Config) (*Parser, error) {
	parser := &Parser{
		cfg:      cfg,
		connPool: newConnectionPool(),
	}
	for dbName, dbCfg := range cfg.Databases {
		initErr := parser.connPool.initConnection(dbName, dbCfg.Name, dbCfg.DSN)
		if initErr != nil {
			return nil, initErr
		}
	}

	return parser, nil
}
