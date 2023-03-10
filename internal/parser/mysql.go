package parser

import (
	"dbseeder/internal/schema"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	getMysqlTablesSql  = "SHOW TABLES"
	getMysqlColumnsSql = "SHOW COLUMNS FROM %s"
)

func MysqlParse(db *sqlx.DB, dbName string) ([]schema.Table, error) {
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
	selectErr := db.Select(&tableNames, getMysqlTablesSql)

	return tableNames, selectErr
}

func getMysqlColumns(db *sqlx.DB, tableName string) ([]string, error) {
	type mysqlColumns struct {
		Field   string `db:"field"`
		Type    string `db:"Type"`
		Null    string `db:"Null"`
		Key     string `db:"Key"`
		Default string `db:"Default"`
		Extra   string `db:"Extra"`
	}

	reqResult := make([]mysqlColumns, 0)
	selectErr := db.Select(&reqResult, fmt.Sprintf(getMysqlColumnsSql, tableName))
	result := make([]string, len(reqResult))
	for _, r := range reqResult {
		result = append(result, r.Field)
	}

	return result, selectErr
}


