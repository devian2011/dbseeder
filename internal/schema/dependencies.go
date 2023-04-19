package schema

import (
	"fmt"
)

type TableDependenceTree struct {
	root     *TableDependenceNode
	tableMap map[string]*TableDependenceNode
}

func NewTree() *TableDependenceTree {
	return &TableDependenceTree{
		tableMap: make(map[string]*TableDependenceNode),
	}
}

func (tree *TableDependenceTree) GetNode(dbName, tableName string) *TableDependenceNode {
	code := TableCode(dbName, tableName)
	if _, exists := tree.tableMap[code]; !exists {
		tree.tableMap[code] = &TableDependenceNode{
			Dependencies: make([]*TableDependenceNode, 0),
			Table:        nil,
			DbName:       dbName,
			Code:         code,
			Dependents:   0,
		}
	}

	return tree.tableMap[code]
}

type TableDependenceNode struct {
	Dependencies []*TableDependenceNode
	Table        *Table
	DbName       string
	Code         string
	Dependents   int
}

func (node *TableDependenceNode) GetDependenceTableCode() string {
	return fmt.Sprintf(node.DbName, node.Table.Name)
}

func (node *TableDependenceNode) addTableNotation(table *Table) {
	node.Table = table
}

func (node *TableDependenceNode) dependsOn(dependent *TableDependenceNode) {
	node.Dependencies = append(node.Dependencies, dependent)
}

func (node *TableDependenceNode) HasDependencies() bool {
	return len(node.Dependencies) > 0
}

func (node *TableDependenceNode) addDependent() {
	node.Dependents++
}

func (node *TableDependenceNode) HasDependents() bool {
	return node.Dependents != 0
}

func BuildDependenciesTree(dbs *Databases) *TableDependenceTree {
	tree := NewTree()
	for _, database := range dbs.Databases {
		for k, table := range database.Tables {
			node := tree.GetNode(database.Name, table.Name)
			node.addTableNotation(&database.Tables[k])
			for _, field := range table.Fields {
				// Set foreign key dependence
				if field.IsFkDependence() {
					relationNode := tree.GetNode(field.Depends.ForeignKey.Db, field.Depends.ForeignKey.Table)
					relationNode.addDependent()
					node.dependsOn(relationNode)
				}
			}
		}
	}

	return tree
}

type TreeWalkFunc func(code string, node *TableDependenceNode) error

func (tree *TableDependenceTree) Walk(fn TreeWalkFunc) error {
	for code, node := range tree.tableMap {
		err := fn(code, node)
		if err != nil {
			return err
		}
	}

	return nil
}
