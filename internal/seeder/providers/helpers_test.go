package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlaceholderString(t *testing.T) {
	actual := createPlaceholderStr(5)
	expected := "(?,?,?,?,?)"

	if actual != expected {
		t.Errorf("Actual '%s' must be eq to Expected '%s'", actual, expected)
	}
}

func TestFillPlaceholderString(t *testing.T) {
	actual := fillPlaceholdersString("(?,?)", 3)
	expected := "(?,?),(?,?),(?,?)"

	assert.Equalf(t, actual, expected, "Actual '%s' must be eq to Expected '%s'", actual, expected)
}

func TestGenerateInsertSql(t *testing.T) {
	actual := generateInsertSQL("info", []string{"id", "name", "val"}, 3)
	expected := "INSERT INTO info (id,name,val) VALUES (?,?,?),(?,?,?),(?,?,?)"

	assert.Equalf(t, actual, expected, "Actual '%s' must be eq to Expected '%s'", actual, expected)
}
