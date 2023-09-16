package cavage

import (
	"fediverse/nullable"
	"fediverse/security"
	"net/http"
)

// AddSignatureToRequest is an opinionated function that adds a signature to an
// HTTP request, ideally to be sent out to a server.
func AddSignatureToRequest(
	req *http.Request,
	params Params,
	signer security.ToStringSigner,
) error {
	ssi := SigningStringInfo{
		Method:          req.Method,
		Path:            req.URL.Path,
		Created:         nullable.Just(params.Created),
		Headers:         req.Header.Clone(),
		ExpectedHeaders: params.Headers.ValueOrDefault([]string{created}),
	}
	signature, err := signer.Sign([]byte(ssi.ConstructSigningString()))
	if err != nil {
		return err
	}
	p := ParamsWithSignature{
		Params:    params,
		Signature: signature,
	}
	req.Header.Set("Signature", p.String())
	return nil
}
