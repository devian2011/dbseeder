package seeder

import "fmt"

type RelationValues struct {
	//TODO: Optimize store values
	Tables map[string][]map[string]any
}

func (t *RelationValues) Add(tableCode string, rows []map[string]any) {
	t.Tables[tableCode] = rows
}

func (t *RelationValues) Get(tableCode string) ([]map[string]any, error) {
	if v, exists := t.Tables[tableCode]; exists {
		return v, nil
	}
	return nil, fmt.Errorf("values for table: %s didn't be generated", tableCode)
}

func (t *RelationValues) isTableGenerated(tableCode string) bool {
	if _, exists := t.Tables[tableCode]; exists {
		return true
	}
	return false
}
