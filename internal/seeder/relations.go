package seeder

import "fmt"

type RelationStore struct {
	Tables map[string][]map[string]any
}

func NewRelationValues() *RelationStore {
	return &RelationStore{Tables: make(map[string][]map[string]any, 1000)}
}

func (t *RelationStore) Add(tableCode string, rows []map[string]any) {
	t.Tables[tableCode] = rows
}

func (t *RelationStore) Get(tableCode string) ([]map[string]any, error) {
	if v, exists := t.Tables[tableCode]; exists {
		return v, nil
	}
	return nil, fmt.Errorf("values for table: %s didn't be generated", tableCode)
}

func (t *RelationStore) IsTableDataGenerated(tableCode string) bool {
	if _, exists := t.Tables[tableCode]; exists {
		return true
	}
	return false
}
