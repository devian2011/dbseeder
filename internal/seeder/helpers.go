package seeder

import (
	"fmt"
	"strings"
)

func createPlaceholderStr(columnsCount int) string {
	byteCnt := columnsCount * 2
	str := make([]byte, byteCnt*2)
	for c := 0; c <= byteCnt; c += 2 {
		str[c] = '?'
		str[c+1] = ','
	}

	return "(" + string(str[:byteCnt-1]) + ")"
}

func fillPlaceholdersString(placeholderString string, vCount int) string {
	var b strings.Builder
	for c := 1; c <= vCount; c++ {
		if c == vCount {
			b.WriteString(placeholderString)
		} else {
			b.WriteString(placeholderString)
			b.WriteString(",")
		}
	}

	return b.String()
}

func generateInsertSQL(tableName string, columns []string, valuesCount int) string {
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(columns, ","),
		fillPlaceholdersString(createPlaceholderStr(len(columns)), valuesCount))
}
