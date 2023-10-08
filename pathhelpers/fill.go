package pathhelpers

import "strings"

func FillFields(path string, params map[string]string) string {
	for k, v := range params {
		path = strings.Replace(path, ":"+k, v, -1)
	}
	return path
}
