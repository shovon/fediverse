package cavage

import (
	"errors"
	"net/http"
)

func DeriveSignatureString(params ParamsWithSignature, header http.Header) (string, error) {
	// expectedHeaders := params.Headers.ValueOrDefault([]string{"(created)"})

	// return SigningStringInfo{}, nil
	return "", errors.New("not implemented")
}
