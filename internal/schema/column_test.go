package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOrderColumns(t *testing.T) {
	fields := map[string]Field{
		"fullname": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends: Dependence{
				Expression: ExpressionDependence{
					Expression: "",
					Rows: []string{
						"firstname", "lastname",
					},
				},
				ForeignKey: ForeignDependence{},
			},
			List: nil,
		},
		"firstname": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends:    Dependence{},
			List:       nil,
		},
		"dt": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends: Dependence{
				Expression: ExpressionDependence{
					Expression: "",
					Rows: []string{
						"fullname",
					},
				},
				ForeignKey: ForeignDependence{},
			},
			List: nil,
		},
		"lastname": {
			Type:       "",
			Generation: "",
			Plugins:    nil,
			Depends:    Dependence{},
			List:       nil,
		},
	}
	expected := []string{"firstname", "lastname", "fullname", "dt"}
	sorter := columnSorter{fields: fields}
	actual, err := sorter.sort()
	if err != nil {
		t.Errorf("test cannot return error. err: %s", err.Error())
	} else {
		assert.Equal(t, expected, actual)
	}
}
