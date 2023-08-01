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

func (schema Schema) Explain() error {
	truncateErr := schema.Tree.WalkAsc(func(node *TreeNode) error {
		fmt.Printf("truncate %s \n", node.Code)
		return nil
	})

	if truncateErr != nil {
		return truncateErr
	}

	generateErr := schema.Tree.WalkDesc(func(node *TreeNode) error {
		fmt.Printf("generate %s\n", node.Code)
		return nil
	})

	return generateErr
}

func (schema Schema) Check() error {
	return schema.Tree.Iterate(
		func(t *Tree, code string, node *TreeNode) error {
			for fieldName, field := range node.Table.Fields {
				if field.IsFkDependence() && field.Depends.ForeignKey.Type == OneToOne {
					depNode := schema.Tree.GetNode(field.Depends.ForeignKey.Db, field.Depends.ForeignKey.Table)
					if node.Table.GetRowsCount() != depNode.Table.GetRowsCount() {
						return fmt.Errorf("in fk oneToOne relation count of rows MUST be equal (%s.%s.%s - %s.%s.%s)",
							depNode.DbName, depNode.Table.Name, field.Depends.ForeignKey.Field,
							node.DbName, node.Table.Name, fieldName)
					}
				}
				if field.IsExpressionDependence() {
					for _, r := range field.Depends.Expression.Rows {
						if node.Table.Fields[r].Generation == GenerationTypeDb {
							return fmt.Errorf("row with expression cannot use db generated values %s.%s.%s - %s))",
								node.DbName, node.Table.Name, fieldName, r)
						}
					}
				}
			}

			return nil
		},
	)
}
