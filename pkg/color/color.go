package color

import "fmt"

type Code string

const (
	Green  Code = "green"
	Yellow Code = "yellow"
	Red    Code = "red"
)

var codeStrMap = map[Code]string{
	Red:    "\u001B[1;31m%s\u001B[0m",
	Green:  "\u001B[1;32m%s\u001B[0m",
	Yellow: "\u001B[1;33m%s\u001B[0m",
}

func ColoredString(code Code, val string) string {
	if tmpl, exists := codeStrMap[code]; exists {
		return fmt.Sprintf(tmpl, val)
	}

	return val
}
