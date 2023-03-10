package parser

import (
	"dbseeder/internal/schema"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

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
	cfg      *schema.Schema
	connPool *connectionPool
}

func NewParser(cfg *schema.Schema) (*Parser, error) {
	parser := &Parser{
		cfg:      cfg,
		connPool: newConnectionPool(),
	}
	for dbName, dbCfg := range cfg.Databases.Databases {
		initErr := parser.connPool.initConnection(dbName, dbCfg.Driver, dbCfg.DSN)
		if initErr != nil {
			return nil, initErr
		}
	}

	return parser, nil
}

func (parser *Parser) Parse() (*schema.Databases, error) {
	databases := make(map[string]*schema.Database, 0)
	for code, conn := range parser.connPool.connections {
		switch parser.cfg.Databases.Databases[code].Driver {
		case "pgx":
			tables, err := PsqlParse(conn, code)
			if err != nil {
				return nil, err
			}
			databases[code] = &schema.Database{
				Driver:     "pgx",
				Name:       code,
				DSN:        parser.cfg.Databases.Databases[code].DSN,
				TablesPath: "",
				Tables:     tables,
			}
		case "mysql":
			tables, err := MysqlParse(conn, code)
			if err != nil {
				return nil, err
			}
			databases[code] = &schema.Database{
				Driver:     "mysql",
				Name:       code,
				DSN:        parser.cfg.Databases.Databases[code].DSN,
				TablesPath: "",
				Tables:     tables,
			}
		default:
			return nil, errors.New("unknown database driver for pare")
		}
	}

	return &schema.Databases{Databases: databases}, nil
}
