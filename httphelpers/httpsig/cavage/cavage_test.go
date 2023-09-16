package cavage

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fediverse/httphelpers"
	"fediverse/security"
	"fediverse/security/rsassapkcsv115sha256"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// func getStringVeifier(key string) security.FromStringVerifier {

// }

func TestSigningVerification(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://localhost", strings.NewReader("hello\n"))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Cool", "hello")
		req.Header.Set("Signature", "headers=\"cool\", signature=\"SNVXJUc7YuaHk9e+Hu6lWQF1TNOOHnCRfdxQX5O+H7qrB4wTZvvbxLVSj7fkzIKmZ97UR0tiNtFGxT1DMIjpTnGfYAuzop2lZQljHmHOm29RvZKe62LoE8M1mjiQciBhRlsVmhmWbcqfIGpHF0m4SW7KmvJDOdxY6K6wPrPYIxTQ+bKEbBVUpVjosAJ2B3kTsLdKqXbU5oml0zczadPG2qriDfzUj5xjY3Oibh4wleoGQlXce0RwXr+0H3AK7wXRuKZjq2wTA5YHQHRDl8QOzTU+JDU42D2529I8XJ7wcOAYgHqz5LBBCGuB2uUq7OUxhmybLZZnOet7qagq+ns39Q==\"")

		publicKeyPEM := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs94w7ycdUeF54NzOc1WQ\n+Oy47dRwagFxdPmyvxnqD2FkAGBF3dcRb6ty0fph6DH5mGa8oV7pRozbYuXg0QYp\nhfdXewT27/IIxRfNJUVGHBgBfyjVy4KQ5S8fvHxaOq5Y5xDrgLDVsf1Xgb8Qz6Cd\nA0xGiUnzH/bbpCmm1H3IvlcWXOAy6fXH2Ghr4curlYAiT7kvsckh+bv0gHAzGAu3\nG6wZA2W68hWIuhSz4jPVv8sIuKMM3OlC3EbbFC6+nwSuI4t/qJhmTzDEeexX816x\nfi+xoAg83tmKZ0w+cywYS/1UVDrtaD77fb8Nrlv7CWwrJ2cl040mKW/OwujqADaO\n8QIDAQAB\n-----END PUBLIC KEY-----\n"

		middleware := VerifySignature(func(r httphelpers.ReadOnlyRequest) (security.FromStringVerifier, error) {
			block, _ := pem.Decode([]byte(publicKeyPEM))
			if block == nil || block.Type != "PUBLIC KEY" {
				// t.Error("failed to decode public key")
				// t.FailNow()
				return nil, errors.New("failed to decode public key")
			}
			publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
			if !ok {
				return nil, errors.New("fatal error")
			}

			verifier := rsassapkcsv115sha256.Base64Verifier(publicKey)

			return verifier, nil
		})

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello\n"))
		}))

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}
	})
}
