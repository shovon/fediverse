package rfc3230

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

type SHA256Digest struct{}

var _ Digester = SHA256Digest{}

func (d SHA256Digest) Token() string {
	return "sha-256"
}

func (d SHA256Digest) Digest(body []byte) (string, error) {
	hash := sha256.Sum256(body)
	return b64.StdEncoding.EncodeToString(hash[:]), nil
}

func TestVerifyDigest(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost", strings.NewReader("hello\n"))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Digest", "sha-256=WJG1tSLV3whtD/CxEPvZ0hu0/HFjrzTQgoai6Eb2vgM=")

	digesters := []Digester{SHA256Digest{}}

	rr := httptest.NewRecorder()

	middleware := VerifyDigest(digesters)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello\n"))
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestVerifyDigestBad(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost", strings.NewReader("hello\n"))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Digest", "sha-256=WJG1tSLV3whtD/CxEPvZ0hu0/HFjryTQgoai6Eb2vgM=")

	digesters := []Digester{SHA256Digest{}}

	rr := httptest.NewRecorder()

	middleware := VerifyDigest(digesters)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello\n"))
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}

func TestVerifyUnknownDigest(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost", strings.NewReader("hello\n"))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Digest", "ha-256=WJG1tSLV3whtD/CxEPvZ0hu0/HFjryTQgoai6Eb2vgM=")

	digesters := []Digester{SHA256Digest{}}

	rr := httptest.NewRecorder()

	middleware := VerifyDigest(digesters)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello\n"))
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	digest := rr.Header().Get("Want-Digest")
	pair, err := ParseWantedDigest(digest)
	if err != nil {
		t.Error("unexpected error")
	}
	if len(pair) != 1 {
		t.Error("expected 1 digest")
		t.FailNow()
	}
	if pair[0].Left != "sha-256" {
		t.Errorf("Expected sha-256, but got %s", pair[0].Left)
	}
	if !pair[0].Right.Equal(decimal.NewFromInt(1)) {
		t.Error("Expected q=1")
	}
}