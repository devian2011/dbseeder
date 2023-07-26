package schema

import (
	"fmt"
)

type Tree struct {
	nodes []*TreeNode
}

type TreeHandleFunc func(tbl *Table, columnOrder []string, dbName, tblCode string) error

func (t *Tree) Walk(fn TreeHandleFunc) error {
	appliedDict := make(map[string]interface{})
	for _, n := range t.nodes {
		if _, exists := appliedDict[n.code]; exists {
			continue
		}
		if err := fn(n.table, n.columnOrder, n.dbName, n.code); err != nil {
			return err
		}
		appliedDict[n.code] = nil
		if err := n.apply(appliedDict, fn); err != nil {
			return err
		}
	}

	return nil
}

type TreeNode struct {
	next        []*TreeNode
	previous    []*TreeNode
	table       *Table
	columnOrder []string
	dbName      string
	code        string
}

func (n *TreeNode) HasDependencies() bool {
	return len(n.next) > 0
}

func (n *TreeNode) HasDependents() bool {
	return len(n.previous) > 0
}

func (n *TreeNode) apply(applied map[string]interface{}, fn TreeHandleFunc) error {
	if _, exists := applied[n.code]; !exists {
		if err := fn(n.table, n.columnOrder, n.dbName, n.code); err != nil {
			return err
		}
	}

	applied[n.code] = nil

	for _, prev := range n.previous {
		if _, exists := applied[prev.code]; exists {
			continue
		}
		if err := prev.apply(applied, fn); err != nil {
			return err
		}
		applied[prev.code] = nil
	}

	for _, next := range n.next {
		if _, exists := applied[next.code]; exists {
			continue
		}

		if err := next.apply(applied, fn); err != nil {
			return err
		}
		applied[next.code] = nil
	}

	return nil
}

func buildTableDependencies(tCode string, tMap map[string]*Table, nMap map[string]*TreeNode) error {
	table := tMap[tCode]
	tNode := &TreeNode{
		next:        make([]*TreeNode, 0),
		previous:    make([]*TreeNode, 0),
		table:       table,
		columnOrder: make([]string, len(table.Fields)),
		dbName:      GetDbFromTableCode(tCode),
		code:        tCode,
	}
	nMap[tCode] = tNode
	for _, fld := range table.Fields {
		if fld.IsFkDependence() {
			dCode := TableCode(fld.Depends.ForeignKey.Db, fld.Depends.ForeignKey.Table)

			if _, exists := nMap[dCode]; !exists {
				if _, exists := tMap[dCode]; !exists {
					return fmt.Errorf("unknown table dependence table %s from %s", tCode, dCode)
				}

				if bErr := buildTableDependencies(dCode, tMap, nMap); bErr != nil {
					return bErr
				}
			}

			nMap[dCode].previous = append(nMap[dCode].previous, tNode)
			tNode.next = append(tNode.next, nMap[dCode])
		}
	}

	return nil
}

func BuildTree(dbs *Databases) (*Tree, error) {
	tableMap := make(map[string]*Table)
	for dbName, database := range dbs.Databases {
		for i, tbl := range database.Tables {
			tableMap[TableCode(dbName, tbl.Name)] = &database.Tables[i]
		}
	}

	tree := &Tree{nodes: make([]*TreeNode, 0, len(tableMap))}
	nodesMap := make(map[string]*TreeNode, len(tableMap))

	for code := range tableMap {
		if _, exists := nodesMap[code]; !exists {
			if err := buildTableDependencies(code, tableMap, nodesMap); err != nil {
				return nil, err
			}
		}
	}

	for _, n := range nodesMap {
		if len(n.previous) == 0 {
			tree.nodes = append(tree.nodes, n)
		}
	}

	return tree, nil
}
