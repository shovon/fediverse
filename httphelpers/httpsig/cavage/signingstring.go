package cavage

import (
	"fediverse/httphelpers"
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
	Created, Expires nullable.Nilable[time.Time]
	Headers          http.Header
	ExpectedHeaders  []string
}

func SigningStringInfoFromRequest(
	params Params,
	requests httphelpers.ReadOnlyRequest,
) SigningStringInfo {
	var ssi SigningStringInfo

	expectedHeaders := params.Headers.ValueOrDefault([]string{"(created)"})

	ssi.Created = nullable.Just(params.Created)
	ssi.Expires = params.Expires
	ssi.Headers = requests.Header.Clone()
	ssi.ExpectedHeaders = expectedHeaders

	return ssi
}

func stringifyNullableTime(nt nullable.Nilable[time.Time]) string {
	return nullable.Then(nt, func(t time.Time) nullable.Nilable[string] {
		return nullable.Just(strconv.FormatInt(t.Unix(), 10))
	}).ValueOrDefault("")
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
				Right: fmt.Sprintf("%s %s", strings.ToLower(ssi.Method), ssi.Path),
			}
		case created:
			return pair.Pair[string, string]{
				Left:  created,
				Right: stringifyNullableTime(ssi.Created),
			}
		case expires:
			return pair.Pair[string, string]{
				Left:  expires,
				Right: stringifyNullableTime(ssi.Expires),
			}
		}

		return pair.Pair[string, string]{Left: s, Right: ssi.Headers.Get(s)}
	})
	return strings.Join(
		slices.Map(
			result,
			func(p pair.Pair[string, string], _ int) string { return strings.ToLower(p.Left) + ": " + p.Right },
		),
		"\n",
	)
}
