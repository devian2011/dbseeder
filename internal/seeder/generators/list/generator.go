package list

import (
	"fmt"
	"math/rand"

	"dbseeder/internal/schema"
)

func Generate(fieldName string, fieldVal schema.Field) (any, error) {
	if len(fieldVal.List) <= 0 {
		return nil, fmt.Errorf("empty list for %s", fieldName)
	}

	rIndex := rand.Intn(len(fieldVal.List) - 1)
	return fieldVal.List[rIndex], nil
}
