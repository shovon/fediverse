package id

import (
	"crypto/rand"
)

func Generate() (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
