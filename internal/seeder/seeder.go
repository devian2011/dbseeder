package seeder

import (
	"errors"
	"fmt"
	"sync"

	gbe "github.com/devian2011/go_basic_extension"

	"dbseeder/internal/modifiers"
	"dbseeder/internal/schema"
	"dbseeder/internal/seeder/generators/dependence"
	"dbseeder/internal/seeder/generators/fake"
	"dbseeder/internal/seeder/generators/list"
)

type Seeder struct {
	db             *Db
	schema         *schema.Schema
	relValues      *RelationStore
	generatorSPool *sync.Pool
}

func NewSeeder(sch *schema.Schema, plugins *modifiers.ModifierStore) (*Seeder, error) {
	relValues := NewRelationValues()
	db, err := NewDb(sch.Databases.Databases)
	if err != nil {
		return nil, err
	}

	return &Seeder{
		db:        db,
		schema:    sch,
		relValues: relValues,
		generatorSPool: &sync.Pool{New: func() any {
			return &tableGenerator{
				relationValues: relValues,
				plugins:        plugins,
			}
		}},
	}, nil
}

func (seeder *Seeder) Run() error {
	// Truncate tables
	err := seeder.schema.Tree.WalkAsc(seeder.truncateTableFn)

	if err != nil {
		seeder.db.pool.rollback()
		return err
	}

	// Fill data to tables
	err = seeder.schema.Tree.WalkDesc(seeder.valueGenerationFn)

	if err != nil {
		seeder.db.pool.rollback()
	} else {
		seeder.db.pool.commit()
	}

	return err
}

func (seeder *Seeder) truncateTableFn(node *schema.TreeNode) error {
	if seeder.relValues.IsTableDataGenerated(node.Code) {
		return nil
	}
	/// If table has dependencies - truncate it
	if node.Table.Action == schema.ActionGenerate {
		truncateErr := seeder.db.truncate(node.DbName, node.Table.Name)
		if truncateErr != nil {
			return truncateErr
		}
	}

	return nil
}

func (seeder *Seeder) valueGenerationFn(node *schema.TreeNode) error {
	switch node.Table.Action {
	case schema.ActionGenerate:
		td := seeder.generatorSPool.Get().(*tableGenerator)
		defer func() {
			td.erase()
			seeder.generatorSPool.Put(td)
		}()

		tableDataErr := td.init(node.Code, node.ColumnOrder, node.Table)
		if tableDataErr != nil {
			return tableDataErr
		}

		tableGenErr := td.generate()
		if tableGenErr != nil {
			return tableGenErr
		}

		insertErr := seeder.db.insert(node.DbName, node.Table.Name, node.ColumnOrder, td.values)
		if insertErr != nil {
			return insertErr
		}

		result, err := seeder.db.loadDataFromDb(node.DbName, node.Table.Name)
		if err == nil {
			seeder.relValues.Add(node.Code, result)
		}
		return err

	case schema.ActionGet:
		result, err := seeder.db.loadDataFromDb(node.DbName, node.Table.Name)
		if err == nil {
			seeder.relValues.Add(node.Code, result)
		}
		return err
	default:
		return fmt.Errorf("unknown action for table: %s", node.Code)
	}
}

type tableGenerator struct {
	rowsCount      int
	columnOrder    []string
	relationValues *RelationStore
	tbl            *schema.Table
	plugins        *modifiers.ModifierStore
	values         []any
	rowsHashes     map[string]bool
	relations      map[string]map[int]bool
}

func (generator *tableGenerator) init(tblCode string, columnOrder []string, tbl *schema.Table) error {
	/// Found length for generated values
	rowsCount := tbl.GetRowsCount()
	if rowsCount == 0 {
		return fmt.Errorf(
			"count rows and fill part for table: %s is empty, set count generated rows or make fill data",
			tblCode)
	}

	generator.tbl = tbl
	generator.columnOrder = columnOrder
	generator.rowsCount = rowsCount

	generator.values = make([]any, 0, rowsCount*len(tbl.Fields))
	generator.rowsHashes = make(map[string]bool, rowsCount)
	generator.relations = make(map[string]map[int]bool, 0)

	return nil
}

func (generator *tableGenerator) erase() {
	generator.rowsCount = 0
	generator.values = nil
	generator.rowsHashes = nil
	generator.relations = nil
	generator.tbl = nil
}

func (generator *tableGenerator) generate() error {
	for c := 0; c < generator.rowsCount; c++ {
		exists := true
		for ok := true; ok; ok = exists {
			rowValues, rowGeneratedErr := generator.generateRow(c)
			if rowGeneratedErr != nil {
				return rowGeneratedErr
			}

			// It table require no duplicates - check that this values has no identical rows
			if generator.tbl.NoDuplicates {
				sliceHash := gbe.HashSlice(rowValues)
				if _, exists = generator.rowsHashes[sliceHash]; !exists {
					generator.values = append(generator.values, rowValues...)
					generator.rowsHashes[sliceHash] = true
				}
			} else {
				generator.values = append(generator.values, rowValues...)
				exists = false
			}
		}
	}

	return nil
}

var ErrGetFromDb = errors.New("get field val from db")

func (generator *tableGenerator) generateRow(rowNumber int) ([]any, error) {
	rowValues := make(map[string]any, len(generator.columnOrder))
	values := make([]any, 0, len(generator.columnOrder))
	for _, fieldName := range generator.columnOrder {
		if len(generator.tbl.Fill) > rowNumber {
			if v, exists := generator.tbl.Fill[rowNumber][fieldName]; exists {
				v, err := generator.plugins.ApplyList(generator.tbl.Fields[fieldName].Plugins, v)
				if err != nil {
					return nil, err
				}
				values = append(values, v)
				rowValues[fieldName] = v
				continue
			}
		}

		v, err := generator.generateFieldValue(fieldName, rowValues)

		if err == ErrGetFromDb {
			continue
		}

		if err != nil {
			return nil, err
		}

		values = append(values, v)
		rowValues[fieldName] = v
	}

	return values, nil
}

func (generator *tableGenerator) generateFieldValue(fieldName string, rowValues map[string]any) (any, error) {
	fieldValue := generator.tbl.Fields[fieldName]
	switch fieldValue.Generation {
	case schema.GenerationTypeConst:
		return generator.plugins.ApplyList(fieldValue.Plugins, fieldValue.ConstValue)
	case schema.GenerationTypeFaker:
		vl, err := fake.Generate(fieldName, fieldValue)
		if err != nil {
			return vl, err
		}
		return generator.plugins.ApplyList(fieldValue.Plugins, vl)
	case schema.GenerationTypeList:
		vl, err := list.Generate(fieldName, fieldValue)
		if err != nil {
			return vl, err
		}
		return generator.plugins.ApplyList(fieldValue.Plugins, vl)
	case schema.GenerationDepends:
		if fieldValue.IsFkDependence() {
			return dependence.GenerateForeign(fieldValue, generator.relationValues, generator.relations)
		}
		if fieldValue.IsExpressionDependence() {
			return dependence.GenerateExpression(fieldValue, rowValues)
		}
		return nil, fmt.Errorf("unknown dependence field generation type for %s", fieldName)
	case schema.GenerationTypeDb:
		return nil, ErrGetFromDb
	default:
		return nil, fmt.Errorf("unknown field generation type for %s", fieldName)
	}
}
