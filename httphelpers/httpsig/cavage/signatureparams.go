package cavage

import "strings"

const (
	KeyID     = "keyId"
	Signature = "signature"
	Algorithm = "algorithm"
	Created   = "created"
	Expires   = "expires"
	Headers   = "headers"
)

func ParseSignatureParams(params string) map[string]string {
	result := make(map[string]string)
	for _, param := range strings.Split(params, ",") {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), "\"")
	}
	return result
}
