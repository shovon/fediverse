package cavage

import (
	"fediverse/functional"
	"fediverse/slices"
	"net/http"
	"strings"
)

// - HTTP Message Signature
// - Signer
// - Verifier
// - HTTP Message Component
// - Derived Component
// - HTTP Message Component Name
// - HTTP Message Component Identifier
// - HTTP Message Component Value
// - Covered Components
// - Signature Base
// - HTTP Message Signature Algorithm
// - Key Material
// - Creation Time
// - Expiration Time
// - Target Message
// - Signature Context

func ComputeSigningMessage(query []string, urlpath string, header http.Header) string {
	clean := functional.Compose(strings.TrimSpace, strings.ToLower)

	return strings.Join(slices.Wrapper[string](query).Map(clean).Map(func(s string) string {
		switch s {
		case "(request-target)":
		case "(created)":
		case "(expires)":
		}
		return s + header.Get(s)
	}), "\n")
}
