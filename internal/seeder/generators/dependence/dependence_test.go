package dependence

import (
	"dbseeder/internal/schema"
	"testing"
)

func TestGenerateExpression(t *testing.T) {
	fld := schema.Field{
		Type:       "",
		Generation: "",
		Plugins:    nil,
		Depends: schema.Dependence{
			Expression: schema.ExpressionDependence{
				Expression: "row.one + ' ' + row.two",
				Rows:       []string{"one", "two"},
			},
			ForeignKey: schema.ForeignDependence{},
		},
		List: nil,
	}

	expected := "hello world"
	actual, err := GenerateExpression(fld, map[string]any{"one": "hello", "two": "world"})
	if err != nil {
		t.Errorf("get generation expression error: %s", err.Error())
	}

	if expected != actual.(string) {
		t.Errorf("values is not actual. Actual - %s, Expected - %s", actual, expected)
	}
}
