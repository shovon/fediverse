package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

type PemPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

func ToPemPair(key *rsa.PrivateKey) PemPair {
	privateKeyPEMBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	privateKeyPEM := pem.EncodeToMemory(privateKeyPEMBlock)

	publicKey := &key.PublicKey
	publicKeyPEMBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPEMBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(publicKeyPEMBlock)

	return PemPair{PrivateKey: privateKeyPEM, PublicKey: publicKeyPEM}
}
