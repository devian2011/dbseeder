package schema

import (
	"fmt"
)

//TODO: Refactor tree walk functions. Change before and after functions

type Tree struct {
	nodes   []*TreeNode
	nodeMap map[string]*TreeNode
}

func (t *Tree) Iterate(fn func(t *Tree, c string, n *TreeNode)) {
	for c, n := range t.nodeMap {
		fn(t, c, n)
	}
}

type TreeHandleFunc func(tbl *Table, columnOrder []string, dbName, tblCode string) error

const (
	nodeStartProcessing = iota
	nodeApplyingPrevious
	nodeAppliedPrevious
	nodeAppliedBeforeNextFn
	nodeApplyingNext
	nodeAppliedNext
	nodeAppliedAfterNextFn
	nodeProcessed
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

func (t *Tree) Walk(beforeFn TreeHandleFunc, afterFn TreeHandleFunc) error {
	m := &executionMap{
		status: make(map[string]int, len(t.nodeMap)),
	}
	for _, n := range t.nodes {
		if err := n.apply(m, beforeFn, afterFn); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tree) GetNode(dbName, tableName string) *TreeNode {
	return t.nodeMap[TableCode(dbName, tableName)]
}

type TreeNode struct {
	table *Table

	dbName string
	code   string

	next     []*TreeNode
	previous []*TreeNode

	columnOrder []string
}

func (n *TreeNode) HasDependencies() bool {
	return len(n.next) > 0
}

func (n *TreeNode) HasDependents() bool {
	return len(n.previous) > 0
}

func (n *TreeNode) apply(m *executionMap, beforeFn TreeHandleFunc, afterFn TreeHandleFunc) error {
	if m.inState(n.code, nodeProcessed) {
		return nil
	}

	if !m.handling(n.code) {
		m.setState(n.code, nodeStartProcessing)
	}

	if m.inState(n.code, nodeStartProcessing) && !m.inState(n.code, nodeApplyingPrevious) {
		m.setState(n.code, nodeApplyingPrevious)
		for _, prev := range n.previous {
			if m.handling(prev.code) {
				continue
			}
			if err := prev.apply(m, beforeFn, afterFn); err != nil {
				return err
			}
		}
		m.setState(n.code, nodeAppliedPrevious)
	}

	if m.inState(n.code, nodeAppliedPrevious) && !m.inState(n.code, nodeAppliedBeforeNextFn) {
		if bErr := beforeFn(n.table, n.columnOrder, n.dbName, n.code); bErr != nil {
			return bErr
		}
		m.setState(n.code, nodeAppliedBeforeNextFn)
	}

	if m.inState(n.code, nodeAppliedBeforeNextFn) && !m.inState(n.code, nodeApplyingNext) {
		m.setState(n.code, nodeApplyingNext)
		for _, next := range n.next {
			if err := next.apply(m, beforeFn, afterFn); err != nil {
				return err
			}
		}
		m.setState(n.code, nodeAppliedNext)
	}

	if m.inState(n.code, nodeAppliedNext) && !m.inState(n.code, nodeAppliedAfterNextFn) {
		if aErr := afterFn(n.table, n.columnOrder, n.dbName, n.code); aErr != nil {
			return aErr
		}
		m.setState(n.code, nodeAppliedAfterNextFn)
	}

	if m.inState(n.code, nodeAppliedAfterNextFn) {
		m.setState(n.code, nodeProcessed)
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

	var err error
	srt := columnSorter{fields: table.Fields}
	tNode.columnOrder, err = srt.sort()
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
