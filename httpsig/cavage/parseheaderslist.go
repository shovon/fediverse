package cavage

import (
	"errors"
	"fediverse/slices"
	"regexp"
	"strings"
)

// TODO: unit test this.

func ParseHeadersList(headersList string) ([]string, error) {
	reg, err := regexp.Compile("[\r\n]")
	if err != nil {
		// Panic, because we are the one that fucked this up.
		panic(err)
	}

	// Edge cases; if the whitespace character is a newline, then this is a
	// mistake. This edge case will likely never be hit, but if it is, then, this
	// is undefined behaviour, according to the spec, and so, just return an
	// error.
	if reg.Match([]byte(headersList)) {
		return nil, errors.New("the headers list")
	}

	return slices.Wrapper[string](strings.Split(headersList, " \t")).Map(strings.TrimSpace), nil
}
