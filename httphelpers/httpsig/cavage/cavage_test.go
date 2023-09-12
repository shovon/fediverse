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

		// privateKeyPem := "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAs94w7ycdUeF54NzOc1WQ+Oy47dRwagFxdPmyvxnqD2FkAGBF\n3dcRb6ty0fph6DH5mGa8oV7pRozbYuXg0QYphfdXewT27/IIxRfNJUVGHBgBfyjV\ny4KQ5S8fvHxaOq5Y5xDrgLDVsf1Xgb8Qz6CdA0xGiUnzH/bbpCmm1H3IvlcWXOAy\n6fXH2Ghr4curlYAiT7kvsckh+bv0gHAzGAu3G6wZA2W68hWIuhSz4jPVv8sIuKMM\n3OlC3EbbFC6+nwSuI4t/qJhmTzDEeexX816xfi+xoAg83tmKZ0w+cywYS/1UVDrt\naD77fb8Nrlv7CWwrJ2cl040mKW/OwujqADaO8QIDAQABAoIBAG3SeZBcEpPfFvqL\n92YGVbkXWKamMmkXLn4cw93Y5ce0UEnGfoJAAb5sMXQx67vJX7uE5yGkgMx5zq4o\n68bUe1/3sKtFUb0Zy+8DZFegX3lh0vAgL8HNm8jDqB3+01zG/TNAanqt/hxqMhbf\nYVVUnOnZlavXwiG/KUanw9w0XPCR100TAO1EduwpOiBaYG3E0n0wnwQmCbQzGhok\nuQFvbO10LMTz4j4IcJRqE17fcDqO5mpGD6XaNbSTVVbYWoeKuy0bQIl/grRVmaOO\nAHEZijtXYiQK3leAnVpbn7VWC0Uy6xjxwR5/WUc8+IbZKkdlvA8KHaDktF666c6L\nMYdRSBECgYEA5XkgTRLbyc35Rc2UX5YVC13Vog7VukLCBdEFCfepKyh+SSkVI51s\nRd8wpKMr5+I9zGshQt/EcNHJQtyTuV+9ph7c++7ujm1be62mfpUSYvU4U/tN1NN5\nUyqi3asdu7GlnjMMAJVQqM9Bs55HZtOJyAbYYCX6nUb1rB0LxRdhqt0CgYEAyKkU\n5xlNKL2hSO+PmRM5Q3JOpdVvFaHLipsiB1pBQFF4+DH7vQ7SEuzWnPj8H/v4sjwq\nnAKEyzDGCqTAJle28kR50Dq3LlV8rYILNW3etRU60AekMqkVWHddvSGPZBy0kHfD\nUszTG+hxpUy87lYW+p4O7lj2RYjLMIPVwslGASUCgYB3Bim78IEyHnUiQKyqG8WU\nnLo3kaxILWJH9A/CCMSlTx6ZwgZl77x+TXMEomIep5nYUuTws/JHdnkHBjRVXZSX\n+sAyqM7x70UZVIvccmQUEg+CeAH51yrB+YZ6pcwJU+6MrPuXvdsVhFMW4I02h5ia\nuPo4fNqOA4VgHbzCIQuWvQKBgEQnAV14f3dt6MCv9PLFz7YztbinZEJRoKMkC4u5\nwN3Kji6mM11EEl+xJzdLbb6jQxWuT76LMHUezLTCviyHsDBax5DM0HihOmZn+8ya\n04BnhIExhzhZq1FPwXvCUURsZ3uF4cZWoQEikq7VAHpmrQlT87hKaOK3EmQY8tpk\nqC0tAoGBAJ26kOPNLkPsJwAKtbklG8uNJCA5zdCDeUBvwNPtGCCdjJgOQ6v0vYe7\nJ1qcfyQJ6umwqURvj/9v/CAwWlZpjLUWitTJxfbv4OKKfHEosaRnVP1lI6FUAuI1\nIKuv2/sFhYQk/68Y2PVGAtW+bBZ7rH8FqvGBFfffJ8dLg4GgXdhg\n-----END RSA PRIVATE KEY-----\n"
		publicKeyPEM := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs94w7ycdUeF54NzOc1WQ\n+Oy47dRwagFxdPmyvxnqD2FkAGBF3dcRb6ty0fph6DH5mGa8oV7pRozbYuXg0QYp\nhfdXewT27/IIxRfNJUVGHBgBfyjVy4KQ5S8fvHxaOq5Y5xDrgLDVsf1Xgb8Qz6Cd\nA0xGiUnzH/bbpCmm1H3IvlcWXOAy6fXH2Ghr4curlYAiT7kvsckh+bv0gHAzGAu3\nG6wZA2W68hWIuhSz4jPVv8sIuKMM3OlC3EbbFC6+nwSuI4t/qJhmTzDEeexX816x\nfi+xoAg83tmKZ0w+cywYS/1UVDrtaD77fb8Nrlv7CWwrJ2cl040mKW/OwujqADaO\n8QIDAQAB\n-----END PUBLIC KEY-----\n"

		// block, _ := pem.Decode([]byte(privateKeyPem))
		// if block == nil || block.Type != "RSA PRIVATE KEY" {
		// 	t.Error("failed to decode private key")
		// 	t.FailNow()
		// }
		// privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

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
