package cavage

import (
	"errors"
	"fediverse/slices"
	"regexp"
	"strings"
)

// ParseHeadersList parses a headers list, as defined in the spec, into a slice
// of strings.
//
// The headers list is just a comma-separated value of HTTP headers and other
// HTTP signature parameters. These are the "keys" to which values to look up,
// and are used to reconstruct the signature string.
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
		return nil, errors.New("unexpected newline character in headers list")
	}

	return slices.Map(strings.Split(headersList, " "), slices.IgnoreIndex(strings.TrimSpace)), nil
}
