package apphttp

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"fediverse/httphelpers/rfc3230"
	"net/http"
)

type SHA256Digest struct{}

var _ rfc3230.Digester = SHA256Digest{}

func (d SHA256Digest) Token() string {
	return "sha-256"
}

func (d SHA256Digest) Digest(body []byte) (string, error) {
	hash := sha256.Sum256(body)
	return b64.URLEncoding.EncodeToString(hash[:]), nil
}

func SendSigned(req *http.Request, key any, body string) (*http.Response, error) {
	rfc3230.AddDigestsToHeaders(req.Header, []byte(body), []rfc3230.Digester{SHA256Digest{}})
	client := &http.Client{}
	return client.Do(req)
}
