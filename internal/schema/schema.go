package schema

import (
	"fmt"
	"strings"
)

func TableCode(dbName, tableName string) string {
	return fmt.Sprintf("%s.%s", dbName, tableName)
}

func GetDbFromTableCode(code string) string {
	return strings.Split(code, ".")[0]
}

type Schema struct {
	Databases *Databases
	Tree      *Tree
}

func NewSchema(dbs *Databases) (*Schema, error) {
	tree, err := BuildTree(dbs)
	if err != nil {
		return nil, err
	}

	sh := &Schema{
		Databases: dbs,
		Tree:      tree,
	}

	return sh, nil
}

func (schema Schema) Check() error {
	return schema.Tree.Walk(func(tbl *Table, columnOrder []string, dbName, dbCode string) error {
		for fieldName, field := range tbl.Fields {
			if field.IsFkDependence() && field.Depends.ForeignKey.Type == OneToOne {
				depNode := schema.Tree.GetNode(field.Depends.ForeignKey.Db, field.Depends.ForeignKey.Table)
				if tbl.GetRowsCount() != depNode.table.GetRowsCount() {
					return fmt.Errorf("in fk oneToOne relation count of rows MUST be equal (%s.%s.%s - %s.%s.%s)",
						depNode.dbName, depNode.table.Name, field.Depends.ForeignKey.Field,
						dbName, tbl.Name, fieldName)
				}
			}
			if field.IsExpressionDependence() {
				for _, r := range field.Depends.Expression.Rows {
					if tbl.Fields[r].Generation == GenerationTypeDb {
						return fmt.Errorf("row with expression cannot use db generated values %s.%s.%s - %s))",
							dbName, tbl.Name, fieldName, r)
					}
				}
			}
		}

		return nil
	})
}
