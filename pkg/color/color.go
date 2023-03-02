package color

import "fmt"

const (
	greenColor  = "\u001B[1;32m%s\u001B[0m"
	yellowColor = "\u001B[1;33m%s\u001B[0m"
	redColor    = "\u001B[1;31m%s\u001B[0m"
)

type Code string

const (
	Green  Code = "green"
	Yellow Code = "yellow"
	Red    Code = "red"
)

func ColoredString(code Code, val string) string {
	switch code {
	case Green:
		return fmt.Sprintf(greenColor, val)
	case Yellow:
		return fmt.Sprintf(yellowColor, val)
	case Red:
		return fmt.Sprintf(redColor, val)
	default:
		return val
	}
}
