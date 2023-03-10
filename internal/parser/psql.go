package parser

import (
	"dbseeder/internal/schema"
	"github.com/jmoiron/sqlx"
)

const (
	getPsqlTablesSql  = "SELECT table_name FROM information_schema.tables WHERE table_schema='public'"
	getPsqlColumnsSql = "SELECT column_name FROM information_schema.columns WHERE table_name=$1;"
	getPsqlFkSql      = `
SELECT
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM
    information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
             ON tc.constraint_name = kcu.constraint_name
                 AND tc.table_schema = kcu.table_schema
        JOIN information_schema.constraint_column_usage AS ccu
             ON ccu.constraint_name = tc.constraint_name
                 AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name=$1`
)

func PsqlParse(db *sqlx.DB, dbName string) ([]schema.Table, error) {
	tableNames, tableNamesErr := getPsqlTables(db)
	if tableNamesErr != nil {
		return nil, tableNamesErr
	}

	tablesResult := make([]schema.Table, 0, len(tableNames))

	for _, tableName := range tableNames {
		columns, colErr := getPsqlColumns(db, tableName)
		if colErr != nil {
			return nil, colErr
		}

		relationColumns, relColsErr := getPsqlDependencies(db, tableName)
		if relColsErr != nil {
			return nil, relColsErr
		}

		fields := make(map[string]schema.Field, len(columns))
		for _, column := range columns {
			depends := schema.Dependence{}

			if v, exist := relationColumns[column]; exist {
				depends.ForeignKey = schema.ForeignDependence{
					Db:    dbName,
					Table: v.tableName,
					Field: v.fieldName,
					Type:  schema.ManyToOne,
				}
			}

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

func getPsqlTables(db *sqlx.DB) ([]string, error) {
	tableNames := make([]string, 0)
	selectErr := db.Select(&tableNames, getPsqlTablesSql)

	return tableNames, selectErr
}

func getPsqlColumns(db *sqlx.DB, tableName string) ([]string, error) {
	result := make([]string, 0)
	selectErr := db.Select(&result, getPsqlColumnsSql, tableName)

	return result, selectErr
}

type psqlFkDependence struct {
	tableName string `db:"foreign_table_name"`
	fieldName string `db:"foreign_column_name"`
}

func getPsqlDependencies(db *sqlx.DB, tableName string) (map[string]psqlFkDependence, error) {
	type outputStr struct {
		ColName   string `db:"column_name"`
		TableName string `db:"foreign_table_name"`
		FieldName string `db:"foreign_column_name"`
	}

	requestResult := make([]outputStr, 0)
	selectErr := db.Select(&requestResult, getPsqlFkSql, tableName)
	if selectErr != nil {
		return nil, selectErr
	}

	result := make(map[string]psqlFkDependence)
	for _, columnData := range requestResult {
		if _, exists := result[columnData.ColName]; exists {
			continue
		}
		result[columnData.ColName] = psqlFkDependence{
			tableName: columnData.TableName,
			fieldName: columnData.FieldName,
		}
	}

	return result, nil
}
