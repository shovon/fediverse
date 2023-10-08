package routeconst

type Unary string

func (u Unary) Replace(value string) string {
	return string(compileRegexp(`:\w+`).ReplaceAll([]byte(u), []byte(value)))
}
