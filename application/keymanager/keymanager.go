package keymanager

import (
	"crypto/rsa"
	"fediverse/application/config"
	"fediverse/security/rsahelpers"
	"os"
	"path"
)

var privateKey *rsa.PrivateKey

func getPrivateKeyFilename() string {
	return path.Join(config.OutputDir(), "privatekey.pem")
}

func generateAndSavePrivateKey() {
	p, err := rsahelpers.GenerateRSPrivateKey(2048)
	if err != nil {
		panic(err)
	}
	privateKey = p
	os.WriteFile(
		getPrivateKeyFilename(),
		[]byte(rsahelpers.PrivateKeyToPKCS1PEMString(privateKey)),
		os.FileMode(0o600),
	)
}

func init() {
	// Note: remember, the act of saving a private key to a file is not part of
	//   any general-purpose library; this is purely a domain-specific,
	//   application-specific purpose, and in the future, we will be saving the
	//   private key on other mediums (perhaps in-memory).
	//
	//   Perhaps we are going to eliminate the requirement to save the private
	//   key to disk in the future.

	if b, err := os.ReadFile(getPrivateKeyFilename()); err != nil {
		generateAndSavePrivateKey()
	} else {
		p, err := rsahelpers.ParsePKCS1PrivateKeyString(string(b))
		if err != nil {
			generateAndSavePrivateKey()
		} else {
			privateKey = p
		}
	}
}

// GetPrivateKey get
func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}
