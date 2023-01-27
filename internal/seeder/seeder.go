package seeder

import (
	"dbseeder/internal/schema"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Seeder struct {
	connections map[string]*sqlx.DB
	trxs        map[string]*sqlx.Tx
	schema      *schema.Schema
}

type TableValues struct {
	Tables map[string][]map[string]any
}

func NewSeeder(sch *schema.Schema) (*Seeder, error) {
	seeder := &Seeder{
		connections: make(map[string]*sqlx.DB),
		schema:      sch,
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

func (seeder *Seeder) Run() {

}
