package seeder

import (
	"dbseeder/internal/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOrderedColumns(t *testing.T) {
	fields := map[string]schema.Field{
		"fullname": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends: schema.Dependence{
				Expression: schema.ExpressionDependence{
					Expression: "",
					Rows: []string{
						"firstname", "lastname",
					},
				},
				ForeignKey: schema.ForeignDependence{},
			},
			List: nil,
		},
		"firstname": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends:    schema.Dependence{},
			List:       nil,
		},
		"dt": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends: schema.Dependence{
				Expression: schema.ExpressionDependence{
					Expression: "",
					Rows: []string{
						"fullname",
					},
				},
				ForeignKey: schema.ForeignDependence{},
			},
			List: nil,
		},
		"lastname": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends:    schema.Dependence{},
			List:       nil,
		},
	}
	expected := []string{"firstname", "lastname", "fullname", "dt"}
	wait := &orderedColumns{cols: make([]string, 0)}

	for f := range fields {
		getOrderedColumns(wait, f, fields)
	}
	/// Check only last points (this points must be right ordered)
	for c := 2; c < 4; c++ {
		assert.Equalf(t, expected[c], wait.cols[c], "wrong columns order Actual: %+v Expected %+v", wait.cols, expected)
	}
}
