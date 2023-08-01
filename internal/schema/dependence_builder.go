package schema

import (
	"fmt"
)

//TODO: Refactor tree walk functions. Change before and after functions

type Tree struct {
	nodes   []*TreeNode
	nodeMap map[string]*TreeNode
}

func (t *Tree) Iterate(fn func(t *Tree, c string, n *TreeNode) error) error {
	for c, n := range t.nodeMap {
		if err := fn(t, c, n); err != nil {
			return err
		}
	}

	return nil
}

type TreeHandleFunc func(node *TreeNode) error

const (
	nodeProcessed = iota
)

type executionMap struct {
	status map[string]int
}

func (e *executionMap) handling(code string) bool {
	_, exists := e.status[code]
	return exists
}

func (e *executionMap) inState(code string, state int) bool {
	if st, exists := e.status[code]; exists {
		return st >= state
	}

	return false
}

func (e *executionMap) setState(code string, state int) {
	e.status[code] = state
}

type WalkHandlers struct {
	BeforeProcessing TreeHandleFunc
	Processing       TreeHandleFunc
	AfterProcessing  TreeHandleFunc
}

func (t *Tree) WalkAsc(handler TreeHandleFunc) error {
	m := &executionMap{
		status: make(map[string]int, len(t.nodeMap)),
	}
	for _, n := range t.nodes {
		if m.handling(n.Code) {
			continue
		}
		if err := n.applyAsc(m, handler); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tree) WalkDesc(handler TreeHandleFunc) error {
	m := &executionMap{
		status: make(map[string]int, len(t.nodeMap)),
	}
	for _, n := range t.nodes {
		if m.handling(n.Code) {
			continue
		}
		if err := n.applyDesc(m, handler); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tree) GetNode(dbName, tableName string) *TreeNode {
	return t.nodeMap[TableCode(dbName, tableName)]
}

type TreeNode struct {
	Table *Table

	DbName string
	Code   string

	next     []*TreeNode
	previous []*TreeNode

	ColumnOrder []string
}

func (n *TreeNode) applyForPrevious(handleFunc TreeHandleFunc) error {
	for _, prev := range n.previous {
		if err := prev.applyForPrevious(handleFunc); err != nil {
			return err
		}
	}
	return handleFunc(n)
}

func (n *TreeNode) applyForNext(handleFunc TreeHandleFunc) error {
	for _, next := range n.next {
		if err := next.applyForNext(handleFunc); err != nil {
			return err
		}
	}

	return handleFunc(n)
}

func (n *TreeNode) HasDependencies() bool {
	return len(n.next) > 0
}

func (n *TreeNode) HasDependents() bool {
	return len(n.previous) > 0
}

func (n *TreeNode) applyAsc(m *executionMap, handler TreeHandleFunc) error {
	m.setState(n.Code, nodeProcessed)

	wrapper := func(node *TreeNode) error {
		if m.handling(node.Code) {
			return nil
		}
		m.setState(node.Code, nodeProcessed)
		return node.applyAsc(m, handler)
	}

	if previousProcessingErr := n.applyForPrevious(wrapper); previousProcessingErr != nil {
		return previousProcessingErr
	}

	if errProcessing := handler(n); errProcessing != nil {
		return errProcessing
	}

	if nextProcessing := n.applyForNext(wrapper); nextProcessing != nil {
		return nextProcessing
	}

	return nil
}

func (n *TreeNode) applyDesc(m *executionMap, handler TreeHandleFunc) error {
	m.setState(n.Code, nodeProcessed)

	wrapper := func(node *TreeNode) error {
		if m.handling(node.Code) {
			return nil
		}
		m.setState(node.Code, nodeProcessed)
		return node.applyDesc(m, handler)
	}

	if nextProcessing := n.applyForNext(wrapper); nextProcessing != nil {
		return nextProcessing
	}

	if errProcessing := handler(n); errProcessing != nil {
		return errProcessing
	}

	if previousProcessingErr := n.applyForPrevious(wrapper); previousProcessingErr != nil {
		return previousProcessingErr
	}

	return nil
}

func buildTableDependencies(tCode string, tMap map[string]*Table, nMap map[string]*TreeNode) error {
	table := tMap[tCode]
	tNode := &TreeNode{
		next:        make([]*TreeNode, 0),
		previous:    make([]*TreeNode, 0),
		Table:       table,
		ColumnOrder: make([]string, len(table.Fields)),
		DbName:      GetDbFromTableCode(tCode),
		Code:        tCode,
	}
	nMap[tCode] = tNode

	var err error
	srt := columnSorter{fields: table.Fields}
	tNode.ColumnOrder, err = srt.sort()
	if err != nil {
		return err
	}

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

	tree := &Tree{
		nodes:   make([]*TreeNode, 0, len(tableMap)),
		nodeMap: make(map[string]*TreeNode, len(tableMap)),
	}

	for code := range tableMap {
		if _, exists := tree.nodeMap[code]; !exists {
			if err := buildTableDependencies(code, tableMap, tree.nodeMap); err != nil {
				return nil, err
			}
		}
	}

	for _, n := range tree.nodeMap {
		if len(n.previous) == 0 {
			tree.nodes = append(tree.nodes, n)
		}
	}

	return tree, nil
}
