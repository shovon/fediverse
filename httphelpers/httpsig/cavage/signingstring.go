package cavage

import (
	"fediverse/nullable"
	"fediverse/pair"
	"fediverse/slices"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SigningStringInfo struct {
	Method, Path     string
	Created, Expires nullable.Nullable[time.Time]
	Headers          http.Header
	ExpectedHeaders  []string
}

func stringifyNullableTime(nt nullable.Nullable[time.Time]) string {
	return fmt.Sprintf("(created): %s", nullable.Then(nt, func(t time.Time) nullable.Nullable[string] {
		return nullable.Just(strconv.FormatInt(t.Unix(), 10))
	}).ValueOrDefault(""))
}

const (
	requestTarget = "(request-target)"
	created       = "(created)"
	expires       = "(expires)"
)

func (ssi SigningStringInfo) ConstructSigningString() string {
	result := slices.Map(ssi.ExpectedHeaders, func(s string, _ int) pair.Pair[string, string] {
		switch s {
		case requestTarget:
			return pair.Pair[string, string]{
				Left:  requestTarget,
				Right: ssi.Method,
			}
		case created:
			return pair.Pair[string, string]{
				Left:  created,
				Right: stringifyNullableTime(ssi.Created),
			}
		case expires:
			return pair.Pair[string, string]{
				Left:  expires,
				Right: stringifyNullableTime(ssi.Created),
			}
		}

		return pair.Pair[string, string]{Left: s, Right: ssi.Headers.Get(s)}
	})
	return strings.Join(
		slices.Map(
			result,
			func(p pair.Pair[string, string], _ int) string { return p.Left + ": " + p.Right },
		),
		"\n",
	)
}
