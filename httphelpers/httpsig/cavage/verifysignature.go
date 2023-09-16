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

			params, err := ParseSignatureParams(signatureHeader)
			if err != nil {
				httperrors.Unauthorized().ServeHTTP(w, r)
				return
			}
			rr, err := httphelpers.ToReadOnlyRequest(r)
			if err != nil {
				return
			}
			sigString := SigningStringInfoFromRequest(params.Params, rr).ConstructSigningString()
			if err != nil {
				httperrors.Unauthorized().ServeHTTP(w, r)
				return
			}

			req, err := httphelpers.ToReadOnlyRequest(r)
			if err != nil {
				httperrors.InternalServerError().ServeHTTP(w, r)
				return
			}
			verifier := getverifier(req)

			if err := verifier.Verify([]byte(sigString), params.Signature); err != nil {
				httperrors.Unauthorized().ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
