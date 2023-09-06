package rfc3230

import (
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/pair"
	"fediverse/possibleerror"
	"fediverse/slices"
	"net/http"

	"github.com/shopspring/decimal"
)

// VerifyDigest handles HTTP requests that have rfc3230 Digest headers. If the
// headers are missing, or the token is not supported, then respond with a 403.
//
// Warning: this processor will buffer the entire request body in memory. So,
// as an added precaution, it may be recommended to have something filter
// requests by Content-Length.
func VerifyDigest(digesters []Digester) hh.Processor {
	numerator := decimal.NewFromInt(1)
	denominator := decimal.NewFromInt(int64(len(digesters)))

	factor := numerator.Div(denominator)

	slices.Map(digesters, func(digester Digester, index int) decimal.Decimal {
		return factor.Mul(decimal.NewFromInt(int64(index + 1)))
	})

	mapping := map[string]func([]byte) (string, error){}
	for _, digester := range digesters {
		mapping[digester.Token()] = digester.Digest
	}

	return hh.ProcessorFunc(func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := hh.ProcessBody(r)
				if err != nil {
					httperrors.InternalServerError().ServeHTTP(w, r)
					return
				}
				slices.Map(r.Header.Values("Digest"), func(str string, _ int) possibleerror.PossibleError[pair.Pair[string, string]] {

				})
				middleware(next).ServeHTTP(w, r)
			})
		}
	})
}
