package schema

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
