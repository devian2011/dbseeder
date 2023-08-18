package schema

import (
	"fmt"

	gbe "github.com/devian2011/go_basic_extension"
)

type columnSorter struct {
	fields map[string]Field
}

func (o *columnSorter) sort() ([]string, error) {
	result := make([]string, 0, len(o.fields))
	for fieldName := range o.fields {
		if o.fields[fieldName].Generation == GenerationTypeDb {
			continue
		}
		err := o.depChainSort(&result, fieldName)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (o *columnSorter) depChainSort(sorted *[]string, fieldName string) error {
	if fieldVal, exists := o.fields[fieldName]; exists {
		if fieldVal.Generation == GenerationTypeDb {
			return nil
		}
		if len(fieldVal.Depends.Expression.Rows) > 0 {
			for _, depRowField := range fieldVal.Depends.Expression.Rows {
				err := o.depChainSort(sorted, depRowField)
				if err != nil {
					return err
				}
			}
		}
		g := gbe.Array[string](*sorted)
		if g.NotContains(fieldName) {
			*sorted = append(*sorted, fieldName)
		}

		return nil
	}

	return fmt.Errorf("unknown dependence column field %s", fieldName)
}
