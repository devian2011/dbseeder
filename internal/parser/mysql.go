package parser

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"dbseeder/internal/schema"
)

const (
	getMysqlTablesSQL  = "SHOW TABLES"
	getMysqlColumnsSQL = "SHOW COLUMNS FROM %s"
)

// MysqlParse parse mysql schema and return it to notation
func MysqlParse(db *sqlx.DB) ([]schema.Table, error) {
	tableNames, tableNamesErr := getMysqlTables(db)
	if tableNamesErr != nil {
		return nil, tableNamesErr
	}

	tablesResult := make([]schema.Table, 0, len(tableNames))

	for _, tableName := range tableNames {
		columns, colErr := getMysqlColumns(db, tableName)
		if colErr != nil {
			return nil, colErr
		}

		fields := make(map[string]schema.Field, len(columns))
		for _, column := range columns {
			depends := schema.Dependence{}

			fields[column] = schema.Field{
				Type:       "string",
				Generation: "faker",
				Plugins:    nil,
				Depends:    depends,
				List:       nil,
			}
		}

		tablesResult = append(tablesResult, schema.Table{
			NoDuplicates: false,
			Count:        10,
			Name:         tableName,
			Action:       schema.ActionGenerate,
			Fields:       fields,
			Fill:         nil,
		})
	}

	return tablesResult, nil
}

func getMysqlTables(db *sqlx.DB) ([]string, error) {
	tableNames := make([]string, 0)
	selectErr := db.Select(&tableNames, getMysqlTablesSQL)

	return tableNames, selectErr
}

func getMysqlColumns(db *sqlx.DB, tableName string) ([]string, error) {
	type mysqlColumns struct {
		Field   string         `db:"Field"`
		Type    string         `db:"Type"`
		Null    sql.NullString `db:"Null"`
		Key     sql.NullString `db:"Key"`
		Default sql.NullString `db:"Default"`
		Extra   sql.NullString `db:"Extra"`
	}

	reqResult := make([]mysqlColumns, 0)
	selectErr := db.Select(&reqResult, fmt.Sprintf(getMysqlColumnsSQL, tableName))
	result := make([]string, 0, len(reqResult))
	for _, r := range reqResult {
		result = append(result, r.Field)
	}

	return result, selectErr
}
