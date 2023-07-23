package dependence

import (
	"fmt"
	"math/rand"

	"github.com/antonmedv/expr"

	"dbseeder/internal/schema"
)

type GeneratedValues interface {
	Get(code string) ([]map[string]any, error)
}

func GenerateForeign(fieldVal schema.Field, genValues GeneratedValues, relations map[string]map[int]bool) (any, error) {
	generatedVls, findErr := genValues.Get(fieldVal.Depends.ForeignKey.GetTableCode())
	if findErr != nil {
		return nil, findErr
	}

	rIndex := 0
	if fieldVal.Depends.ForeignKey.Type == schema.OneToOne {
		if _, exists := relations[fieldVal.Depends.ForeignKey.Field]; !exists {
			relations[fieldVal.Depends.ForeignKey.Field] = make(map[int]bool, 0)
			relations[fieldVal.Depends.ForeignKey.Field][rIndex] = true
		} else {
			// Find biggest index from relation map
			for k := range relations[fieldVal.Depends.ForeignKey.Field] {
				if rIndex < k {
					rIndex = k
				}
			}
			rIndex++
		}
		relations[fieldVal.Depends.ForeignKey.Field][rIndex] = true
	} else {
		rIndex = rand.Intn(len(generatedVls))
	}

	return generatedVls[rIndex][fieldVal.Depends.ForeignKey.Field], nil
}

func GenerateExpression(fieldVal schema.Field, rowsVal map[string]any) (any, error) {
	ctx := make(map[string]map[string]any, 1)
	ctx["row"] = make(map[string]any, len(fieldVal.Depends.Expression.Rows))
	for _, row := range fieldVal.Depends.Expression.Rows {
		if val, exists := rowsVal[row]; exists {
			ctx["row"][row] = val
		} else {
			return nil, fmt.Errorf("cannot find ctx field: %s in generated values", row)
		}
	}

	return expr.Eval(fieldVal.Depends.Expression.Expression, ctx)
}
