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

func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}
