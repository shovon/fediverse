package rsahelpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateRSPrivateKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func PrivateKeyToPKCS1PEMBlock(key *rsa.PrivateKey) *pem.Block {
	return &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
}

func PrivateKeyToPKCS1PEMString(key *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(PrivateKeyToPKCS1PEMBlock(key)))
}

func PublicKeyToPKIXBlock(key *rsa.PublicKey) (*pem.Block, error) {
	publicKeyPEMBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPEMBytes,
	}, nil
}

func PublicKeyToPKIXString(key *rsa.PublicKey) (string, error) {
	publicKeyPEMBlock, err := PublicKeyToPKIXBlock(key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(publicKeyPEMBlock)), nil
}
