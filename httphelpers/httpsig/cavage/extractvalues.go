package cavage

import (
	"fediverse/httphelpers/httperrors"
	"net/http"
	"strings"
)

func DeriveSignatureString(header http.Header) (string, error) {
	signatureHeader := header.Get("signature")
	if strings.TrimSpace(header.Get("Signature")) == "" {
		return "", httperrors.Unauthorized()
	}

	params := ParseSignatureParams(signatureHeader)
	params.Headers.HasValue()

	// signature := req.Header.Get("Signature")
	// params := ParseSignatureParams(signature)
	// value := params.Headers.ValueOrDefault("")
	// headersList, err := ParseHeadersList(value)
	// if err != nil {
	// 	return SigningStringInfo{}, err
	// }
	return SigningStringInfo{}, nil
}
