package cavage

import (
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/security"
	"net/http"
	"strings"
)

func VerifySignature(getverifier func(httphelpers.ReadOnlyRequest) security.FromStringVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signatureHeader := r.Header.Get("signature")
			if strings.TrimSpace(r.Header.Get("Signature")) == "" {
				httperrors.Unauthorized().ServeHTTP(w, r)
				return
			}

			params := ParseSignatureParams(signatureHeader)

			req, err := httphelpers.ToReadOnlyRequest(r)
			if err != nil {
				httperrors.InternalServerError().ServeHTTP(w, r)
				return
			}
			verifier := getverifier(req)

			if err := verifier.Verify([]byte(params.String()), params.Signature); err != nil {
				httperrors.Unauthorized().ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
