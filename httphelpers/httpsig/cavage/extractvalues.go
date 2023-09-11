package cavage

import (
	"fediverse/httphelpers"
)

func DeriveSignatureString(req httphelpers.ReadOnlyRequest) (SigningStringInfo, error) {
	// signature := req.Header.Get("Signature")
	// params := ParseSignatureParams(signature)
	// value := params.Headers.ValueOrDefault("")
	// headersList, err := ParseHeadersList(value)
	// if err != nil {
	// 	return SigningStringInfo{}, err
	// }
	return SigningStringInfo{}, nil
}
