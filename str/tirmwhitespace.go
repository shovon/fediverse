package str

import "strings"

func TrimWhitespace(str string) string {
	return strings.Trim(str, " \t\n\r")
}
