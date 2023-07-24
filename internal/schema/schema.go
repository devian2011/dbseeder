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
	Tree      *TableDependenceTree
}

func NewSchema(dbs *Databases) (*Schema, error) {
	tree, err := BuildDependenciesTree(dbs)
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
	return schema.Tree.Walk(func(code string, node *TableDependenceNode) error {
		if node.HasDependencies() {
			for _, dep := range node.Dependencies {
				for fieldName, field := range dep.Table.Fields {
					if field.IsFkDependence() && field.Depends.ForeignKey.Type == OneToOne {
						depNode := schema.Tree.GetNode(field.Depends.ForeignKey.Db, field.Depends.ForeignKey.Table)
						if node.Table.GetRowsCount() >= depNode.Table.GetRowsCount() {
							return fmt.Errorf("in fk oneToOne relation count of rows MUST be equal (%s.%s.%s - %s.%s.%s)",
								depNode.DbName, depNode.Table.Name, field.Depends.ForeignKey.Field,
								dep.DbName, dep.Table.Name, fieldName)
						}
					}
					if field.IsExpressionDependence() {
						for _, r := range field.Depends.Expression.Rows {
							if dep.Table.Fields[r].Generation == GenerationTypeDb {
								return fmt.Errorf("row with expression cannot use db generated values %s.%s.%s - %s))",
									dep.DbName, dep.Table.Name, fieldName, r)
							}
						}
					}
				}
			}
		}

		return nil
	})
}
