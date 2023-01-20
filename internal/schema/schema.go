package schema

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

type Schema struct {
	databases    map[string]Database
	dependencies map[string]string
}

func NewSchema(dbs map[string]Database) *Schema {
	sh := &Schema{
		databases:    dbs,
		dependencies: make(map[string]string),
	}
	sh.build()

	return sh
}

func (s *Schema) buildDependenciesMap() {
	for _, db := range s.databases {
		for _, table := range db.Tables {
			for _, field := range table.Fields {
				for _, depends := range field.Depends.Foreign {
					s.dependencies[tableKey(db.Name, table.Name, field.Name)] = tableKey(depends.Db, depends.Table, depends.Field)
				}
			}
		}
	}
}

func (s *Schema) build() {
	s.buildDependenciesMap()
}

func tableKey(db string, table string, field string) string {
	return fmt.Sprintf("%s.%s.%s", db, table, field)
}

func (s *Schema) Export(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(s.databases)
}

func (s *Schema) GetDependencies() map[string]string {
	return s.dependencies
}
