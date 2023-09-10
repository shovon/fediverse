package rfc3230

import (
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/pair"
	"fediverse/possibleerror"
	"fediverse/slices"
	"fmt"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

func assertNoError[T any](e possibleerror.PossibleError[T]) T {
	v, err := e.Value()
	if err != nil {
		panic(err)
	}
	return v
}

func calculateTokensAndQValues(digesters []Digester) []pair.Pair[string, decimal.Decimal] {
	numerator := decimal.NewFromInt(1)
	denominator := decimal.NewFromInt(int64(len(digesters)))

	factor := numerator.Div(denominator)

	return slices.Map(digesters, func(digester Digester, index int) pair.Pair[string, decimal.Decimal] {
		fmt.Println(digester.Token())
		return pair.Pair[string, decimal.Decimal]{
			Left:  digester.Token(),
			Right: factor.Mul(decimal.NewFromInt(int64(index + 1))),
		}
	})
}

func deriveWantedDigest(tokenQValuePair pair.Pair[string, decimal.Decimal]) string {
	token, qValue := tokenQValuePair.Left, tokenQValuePair.Right
	if qValue.Equal(decimal.NewFromInt(1)) {
		return token
	}
	return token + ";q=" + qValue.String()
}

func DeriveWantDigests(tokenQValuePairs []pair.Pair[string, decimal.Decimal]) string {
	return strings.Join(slices.Map(tokenQValuePairs, slices.IgnoreIndex(deriveWantedDigest)), ", ")
}

// VerifyDigest produces a middleware to handle HTTP requests that have rfc3230
// Digest headers. If the headers are missing, or the token is not supported,
// then respond with a 403.
//
// Warning: this processor will buffer the entire request body in memory. So,
// as an added precaution, it may be recommended to have something filter
// requests by Content-Length.
func VerifyDigest(digesters []Digester) func(http.Handler) http.Handler {
	pairs := calculateTokensAndQValues(digesters)

	mapping := map[string]func([]byte) (string, error){}
	for _, digester := range digesters {
		mapping[digester.Token()] = digester.Digest
	}

	// TODO: yeah this will really need to be thoroughly unit tested.

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := hh.ProcessBody(r)
			if err != nil {
				httperrors.InternalServerError().ServeHTTP(w, r)
				return
			}

			clientDigests := []possibleerror.PossibleError[[]pair.Pair[string, string]]{}
			for _, digestValue := range r.Header.Values("Digest") {
				clientDigests = append(clientDigests, possibleerror.New(ParseDigest(digestValue)))
			}

			for _, pair := range clientDigests {
				_, err := pair.Value()
				if err != nil {
					httperrors.BadRequest().ServeHTTP(w, r)
					return
				}
			}

			noErrors := slices.Map(clientDigests, slices.IgnoreIndex(assertNoError[[]pair.Pair[string, string]]))
			allDigests := slices.Join(noErrors...)

			for _, digest := range allDigests {
				digestFn, ok := mapping[digest.Left]
				if !ok {
					p := append(pairs, pair.Pair[string, decimal.Decimal]{
						Left:  digest.Left,
						Right: decimal.NewFromInt(0),
					})
					r.Header.Add("Want-Digest", DeriveWantDigests(p))
					httperrors.Unauthorized().ServeHTTP(w, r)
					return
				}
				result, err := digestFn(body)
				if err != nil || digest.Right != result {
					httperrors.Unauthorized().ServeHTTP(w, r)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
