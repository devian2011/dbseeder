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

func (tree *TableDependenceTree) GetNode(dbName string, tableName string) *TableDependenceNode {
	code := fmt.Sprintf("%s.%s", dbName, tableName)
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

func (node *TableDependenceNode) AddTableNotation(table *Table) {
	node.Table = table
}

func (node *TableDependenceNode) DependsOn(dependent *TableDependenceNode) {
	node.Dependencies = append(node.Dependencies, dependent)
}

func (node *TableDependenceNode) HasDependencies() bool {
	return len(node.Dependencies) > 0
}

func (node *TableDependenceNode) AddDependent() {
	node.Dependents++
}

func (node *TableDependenceNode) HasDependents() bool {
	return node.Dependents != 0
}

func BuildDependenciesTree(dbs *Databases) *TableDependenceTree {
	tree := NewTree()
	for _, database := range dbs.Databases {
		for _, table := range database.Tables {
			node := tree.GetNode(database.Name, table.Name)
			node.AddTableNotation(&table)
			for _, field := range table.Fields {
				for _, relation := range field.Depends.Foreign {
					relationNode := tree.GetNode(relation.Db, relation.Table)
					relationNode.AddDependent()
					node.DependsOn(relationNode)
				}
			}
		}
	}

	return tree
}

type TreeWalkFunc func(code string, node *TableDependenceNode)

func (tree *TableDependenceTree) Walk(fn TreeWalkFunc) {
	for code, node := range tree.tableMap {
		fn(code, node)
	}
}
