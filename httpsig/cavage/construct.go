package cavage

import (
	"fediverse/nullable"
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

type pair[K any, V any] struct {
	key   K
	value V
}

func (ssi SigningStringInfo) ConstringSigningString() string {
	result := slices.Map(ssi.ExpectedHeaders, func(s string) pair[string, string] {
		switch s {
		case requestTarget:
			return pair[string, string]{
				key:   requestTarget,
				value: ssi.Method,
			}
		case created:
			return pair[string, string]{
				key:   created,
				value: stringifyNullableTime(ssi.Created),
			}
		case expires:
			return pair[string, string]{
				key:   expires,
				value: stringifyNullableTime(ssi.Created),
			}
		}

		return pair[string, string]{key: s, value: ssi.Headers.Get(s)}
	})
	return strings.Join(
		slices.Map(
			result,
			func(p pair[string, string]) string { return p.key + ": " + p.value },
		),
		"\n",
	)
}
