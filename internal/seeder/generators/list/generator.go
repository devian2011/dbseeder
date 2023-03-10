package list

import (
	"dbseeder/internal/schema"
	"errors"
	"fmt"
	"math/rand"
)

func Generate(fieldName string, fieldVal schema.Field) (any, error) {
	if len(fieldVal.List) <= 0 {
		return nil, errors.New(fmt.Sprintf("empty list for %s", fieldName))
	}

	rIndex := rand.Intn(len(fieldVal.List) - 1)
	return fieldVal.List[rIndex], nil
}
