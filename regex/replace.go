package main

import "regexp"

type Source string
type Substitute string

func ReplaceOnce(r *regexp.Regexp, src Source, sub Substitute) string {
	i := r.FindIndex([]byte(src))
	return string(src[:i[0]]) + string(sub) + string(src[i[1]:])
}
