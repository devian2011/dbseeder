package schema

import (
	"errors"
	"fmt"
)

type Schema struct {
	Databases *Databases
	Tree      *TableDependenceTree
}

func NewSchema(dbs *Databases) *Schema {
	sh := &Schema{
		Databases: dbs,
		Tree:      BuildDependenciesTree(dbs),
	}

	return sh
}

func (schema Schema) Check() error {
	return schema.Tree.Walk(func(code string, node *TableDependenceNode) error {
		if node.HasDependencies() {
			for _, dep := range node.Dependencies {
				for fieldName, field := range dep.Table.Fields {
					if field.IsFkDependence() && field.Depends.ForeignKey.Type == OneToOne {
						depNode := schema.Tree.GetNode(field.Depends.ForeignKey.Db, field.Depends.ForeignKey.Table)
						if node.Table.GetRowsCount() >= depNode.Table.GetRowsCount() {
							return errors.New(
								fmt.Sprintf("in fk oneToOne relation count of rows MUST be equal (%s.%s.%s - %s.%s.%s)",
									depNode.DbName, depNode.Table.Name, field.Depends.ForeignKey.Field,
									dep.DbName, dep.Table.Name, fieldName))
						}
					}
					if field.IsExpressionDependence() {
						for _, r := range field.Depends.Expression.Rows {
							if dep.Table.Fields[r].Generation == GenerationTypeDb {
								return errors.New(
									fmt.Sprintf("row with expression cannot use db generated values %s.%s.%s - %s))",
										dep.DbName, dep.Table.Name, fieldName, r))
							}
						}
					}
				}
			}
		}

		return nil
	})
}
