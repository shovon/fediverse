package routeconst

import "regexp"

type Binary string

type source string
type substitute string

func replaceOnce(r *regexp.Regexp, src source, sub substitute) string {
	i := r.FindIndex([]byte(src))
	return string(src[:i[0]]) + string(sub) + string(src[i[1]:])
}

func (u Binary) Replace(first string, second string) string {
	reg := compileRegexp(`:\w+`)
	return replaceOnce(
		reg,
		source(replaceOnce(reg, source(u), substitute(first))),
		substitute(second),
	)
}
