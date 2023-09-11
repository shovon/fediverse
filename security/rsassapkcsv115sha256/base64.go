package rsassapkcsv115sha256

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	b64 "encoding/base64"
	"fediverse/security"
)

// Pakcage rsassapkcsv115sha256 holds an implementation of

type base64Signer struct {
	privateKey *rsa.PrivateKey
}

func Base64Signer(privateKey *rsa.PrivateKey) security.ToStringSigner {
	return base64Signer{privateKey: privateKey}
}

var _ security.ToStringSigner = base64Signer{}

func (b base64Signer) Sign(payload []byte) (string, error) {
	hash := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(nil, b.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(signature), nil
}

type base64Verifier struct {
	publicKey *rsa.PublicKey
}

var _ security.FromStringVerifier = base64Verifier{}

func Base64Verifier(publicKey *rsa.PublicKey) security.FromStringVerifier {
	return base64Verifier{publicKey: publicKey}
}

func (b base64Verifier) Verify(payload []byte, signature string) error {
	hash := sha256.Sum256(payload)

	sig, err := b64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(b.publicKey, crypto.SHA256, hash[:], sig)
}
