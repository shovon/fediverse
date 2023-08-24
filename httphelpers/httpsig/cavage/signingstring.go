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

func (ssi SigningStringInfo) ConstringSigningString() string {
	result := slices.Map(ssi.ExpectedHeaders, func(s string) pair.Pair[string, string] {
		switch s {
		case requestTarget:
			return pair.Pair[string, string]{
				Key:   requestTarget,
				Value: ssi.Method,
			}
		case created:
			return pair.Pair[string, string]{
				Key:   created,
				Value: stringifyNullableTime(ssi.Created),
			}
		case expires:
			return pair.Pair[string, string]{
				Key:   expires,
				Value: stringifyNullableTime(ssi.Created),
			}
		}

		return pair.Pair[string, string]{Key: s, Value: ssi.Headers.Get(s)}
	})
	return strings.Join(
		slices.Map(
			result,
			func(p pair.Pair[string, string]) string { return p.Key + ": " + p.Value },
		),
		"\n",
	)
}
