package exporter

import (
	"dbseeder/internal/schema"
	"fmt"
	"io"
	"strings"
)

// TableDependenceTreeExporter tree exporter
type TableDependenceTreeExporter struct {
	wr   io.Writer
	tree *schema.TableDependenceTree
}

// NewTableDependenceTreeExporter create new dependence tree exporter instance
func NewTableDependenceTreeExporter(wr io.Writer, tree *schema.TableDependenceTree) *TableDependenceTreeExporter {
	return &TableDependenceTreeExporter{
		wr:   wr,
		tree: tree,
	}
}

// Export export table dependence to stdout
func (exporter *TableDependenceTreeExporter) Export() error {
	return exporter.tree.Walk(exporter.nodeHandler)
}

func (exporter *TableDependenceTreeExporter) nodeHandler(code string, node *schema.TableDependenceNode) error {
	fmt.Println(code)
	for _, d := range node.Dependencies {
		fmt.Print(strings.Repeat(" --> ", 1))
		exporter.nodeHandlerLvl(d.Code, d, 2)
		fmt.Println()
	}

	return nil
}

func (exporter *TableDependenceTreeExporter) nodeHandlerLvl(code string, node *schema.TableDependenceNode, lvl int) {
	fmt.Printf("%s ", code)
	for _, d := range node.Dependencies {
		fmt.Print(" --> ")
		exporter.nodeHandlerLvl(d.Code, d, lvl+1)
	}
}
