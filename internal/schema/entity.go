package schema

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Action string
type Generation string
type ForeignKeyType string

const (
	ActionGenerate Action = "generate" // Generate list of values for insert to db
	ActionGet      Action = "get"      // Get list of values from db

	GenerationTypeDb    Generation = "db"      // Data will be generated by DB for example auto_increment and so on
	GenerationTypeFaker Generation = "faker"   // Generate fake data
	GenerationTypeList  Generation = "list"    // Get data from defined list (property List of Field struct)
	GenerationDepends   Generation = "depends" // Mark field that it depends on another field (on same ot other table)

	ManyToOne ForeignKeyType = "manyToOne"
	OneToOne  ForeignKeyType = "oneToOne"
)

type Databases struct {
	Databases map[string]*Database `json:"databases" yaml:"databases"`
}

type Database struct {
	Driver     string  `json:"driver" yaml:"driver"`
	Name       string  `json:"name" yaml:"name"` // Database name
	DSN        string  `json:"dsn" yaml:"dsn"`   // Database DSN
	TablesPath string  `json:"tables_path" yaml:"tablesPath"`
	Tables     []Table `json:"tables" yaml:"tables"`
}

type Table struct {
	/// TODO: Add unique values block
	NoDuplicates bool             `json:"no_duplicates" yaml:"noDuplicates"` // Allow duplicates or no
	Count        int              `json:"count" yaml:"count"`                // Count rows for generate values
	Name         string           `json:"name" yaml:"name,omitempty"`        // Table name
	Action       Action           `json:"action" yaml:"action,omitempty"`    // Action (get from db or generate fake data)
	Fields       map[string]Field `json:"fields" yaml:"fields,omitempty"`
	Fill         []map[string]any `json:"fill" yaml:"fill,omitempty"`
}

func (t *Table) GetRowsCount() int {
	if t.Count <= 0 {
		if len(t.Fill) <= 0 {
			return 0
		} else {
			return len(t.Fill)
		}
	}

	return t.Count
}

func (t *Table) IsLoadFromDb() bool {
	return t.Action == ActionGet
}

type Field struct {
	Type       string     `json:"type" yaml:"type,omitempty"`             // Type of field - string, email, hash, mac, ip and so on... See types constants
	Generation Generation `json:"generation" yaml:"generation,omitempty"` // Generation strategy
	Plugins    []string   `json:"plugins" yaml:"plugins,omitempty"`       // Plugins list for apply
	Depends    Dependence `json:"depends" yaml:"depends,omitempty"`
	List       []any      `json:"list" yaml:"list,omitempty"`
}

func (fld *Field) IsFkDependence() bool {
	return fld.Depends.ForeignKey.Db != "" && fld.Depends.ForeignKey.Table != "" && fld.Depends.ForeignKey.Field != ""
}

func (fld *Field) IsExpressionDependence() bool {
	return fld.Depends.Expression.Expression != ""
}

type Dependence struct {
	Expression ExpressionDependence `json:"expression" yaml:"expression,omitempty"`
	ForeignKey ForeignDependence    `json:"foreign_key" yaml:"foreign,omitempty"`
}

type ForeignDependence struct {
	Db    string         `json:"db" yaml:"db,omitempty"`
	Table string         `json:"table" yaml:"table,omitempty"`
	Field string         `json:"field" yaml:"field,omitempty"`
	Type  ForeignKeyType `json:"type" yaml:"type,omitempty"`
}

func (fd *ForeignDependence) GetTableCode() string {
	return TableCode(fd.Db, fd.Table)
}

type ExpressionDependence struct {
	Expression string   `yaml:"expression" json:"expression,omitempty"`
	Rows       []string `json:"rows" yaml:"rows,omitempty"`
}

type ExpressionForeignField struct {
	Db    string `json:"db" yaml:"db,omitempty"`
	Table string `json:"table" yaml:"table,omitempty"`
	Field string `json:"field" yaml:"field,omitempty"`
}

func NewDatabasesSchemaNotation(mainConf string) (*Databases, error) {
	cfg := &Databases{}
	dbData, readFileErr := os.ReadFile(mainConf)
	if readFileErr != nil {
		return nil, readFileErr
	}
	unmarshalErr := yaml.Unmarshal(dbData, cfg)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	for id, databaseCfg := range cfg.Databases {
		if databaseCfg.TablesPath != "" {
			databaseCfg.TablesPath = strings.ReplaceAll(databaseCfg.TablesPath, "$PWD", filepath.Dir(mainConf))

			walkErr := filepath.Walk(databaseCfg.TablesPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				tableData, tableFileReadErr := os.ReadFile(path)
				if tableFileReadErr != nil {
					return tableFileReadErr
				}

				fileTables := make([]Table, 0, 0)
				tableUnmarshalErr := yaml.Unmarshal(tableData, &fileTables)
				if tableUnmarshalErr != nil {
					return tableUnmarshalErr
				}
				databaseCfg.Tables = append(databaseCfg.Tables, fileTables...)

				return nil
			})
			if walkErr != nil {
				return nil, walkErr
			}
			cfg.Databases[id] = databaseCfg
		}
	}

	return cfg, nil
}
