package cavage

import (
	"fediverse/security"
	"io"
	"net/http"
	"strings"
)

// AddSignatureToRequest is an opinionated function that adds a signature to an
// HTTP request, ideally to be sent out to a server.
func AddSignatureToRequest(req *http.Request, params Params, signer security.ToStringSigner) error {
	slice, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(strings.NewReader(string(slice)))
	signature, err := signer.Sign(slice)
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
