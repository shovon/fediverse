package cavage

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fediverse/httphelpers"
	"fediverse/security"
	"fediverse/security/rsassapkcsv115sha256"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSigningVerification(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://localhost", strings.NewReader("hello\n"))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Cool", "hello")
		req.Header.Set("Signature", "headers=\"cool\" signature=\"CqH5cyhaTEARs8y/ueIIae0sT24FugeKkXliLQaI1uKeJDXoH16fHolLdya41oTFzK3+uAHN++rVEHIHFXGXlTYEm0AEaqsZzqOSCEcVW7IFzlYKXlTjR36q+vXozpM76aaxN3phljNGi8stV8gAM3TOPVM6lXBFP2isfgVzRD5gEtBxMLsySpPqMGeYZqtUhaiZxCWW5EOi7KtVf9Yp5f+MJ9mJqzNc9j7oerxP31dpQhS0XDu5VDtGHqIOsHxtylWqK52RBNjkfVSm4ajPQugeUd9OyNOBlTrvZZFcPiUlNnQduVk9SneqavXP//mwJyP8d7N0gtcUbFw/YeXTiw==\"")

		publicKeyPEM := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs94w7ycdUeF54NzOc1WQ\n+Oy47dRwagFxdPmyvxnqD2FkAGBF3dcRb6ty0fph6DH5mGa8oV7pRozbYuXg0QYp\nhfdXewT27/IIxRfNJUVGHBgBfyjVy4KQ5S8fvHxaOq5Y5xDrgLDVsf1Xgb8Qz6Cd\nA0xGiUnzH/bbpCmm1H3IvlcWXOAy6fXH2Ghr4curlYAiT7kvsckh+bv0gHAzGAu3\nG6wZA2W68hWIuhSz4jPVv8sIuKMM3OlC3EbbFC6+nwSuI4t/qJhmTzDEeexX816x\nfi+xoAg83tmKZ0w+cywYS/1UVDrtaD77fb8Nrlv7CWwrJ2cl040mKW/OwujqADaO\n8QIDAQAB\n-----END PUBLIC KEY-----\n"

		middleware := VerifySignature(func(r httphelpers.ReadOnlyRequest) security.FromStringVerifier {
			block, _ := pem.Decode([]byte(publicKeyPEM))
			if block == nil || block.Type != "PUBLIC KEY" {
				t.Error("failed to decode public key")
				t.FailNow()
			}
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				t.Error(err.Error())
				t.FailNow()
			}
			publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
			if !ok {
				t.Error("failed to parse public key")
				t.FailNow()
			}

			verifier := rsassapkcsv115sha256.Base64Verifier(publicKey)

			return verifier
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
