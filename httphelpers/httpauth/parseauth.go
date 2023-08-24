package httpauth

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Authorization struct {
	Scheme string
	Params string
}

func ParseAuthorization(authorization string) Authorization {
	pattern := regexp.MustCompile("[ \t]+")
	parts := pattern.Split(authorization, -1)
	if len(parts) < 1 {
		return Authorization{}
	}
	scheme := strings.ToLower(parts[0])

	params := string(([]rune(authorization))[:utf8.RuneCountInString(scheme)+1])

	return Authorization{Scheme: scheme, Params: params}
}
