package rsassapkcsv115sha256

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	b64 "encoding/base64"
	"fediverse/security"
)

// Pakcage rsassapkcsv115sha256 holds an implementation of

type base64 struct {
	privateKey *rsa.PrivateKey
}

func Base64(privateKey *rsa.PrivateKey) security.ToStringSigner {
	return base64{privateKey: privateKey}
}

var _ security.ToStringSigner = base64{}

func (b base64) Sign(payload []byte) (string, error) {
	hash := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(nil, b.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(signature), nil
}
