package color

import "fmt"

const (
	greenColor  = "\u001B[1;32m%s\u001B[0m"
	yellowColor = "\u001B[1;33m%s\u001B[0m"
)

type Code string

const (
	Green  Code = "green"
	Yellow Code = "yellow"
)

func ColoredString(code Code, val string) string {
	switch code {
	case Green:
		return fmt.Sprintf(greenColor, val)
	case Yellow:
		return fmt.Sprintf(yellowColor, val)
	default:
		return val
	}
}
