package cavage

import (
	"fediverse/nullable"
	"fmt"
	"strings"
)

// Section 2.1
//
// Considering that the specification does not formally define what the specific
// format the string is, and instead resorts to examples, I decided to just hard
// code the parsing of the signature parameters, rather than to write a more
// generic parser that handles all the gnarly edge cases.

const (
	KeyID     = "keyId"
	Signature = "signature"
	Algorithm = "algorithm"
	Created   = "created"
	Expires   = "expires"
	Headers   = "headers"
)

type SignatureParams struct {
	KeyID     nullable.Nullable[string]
	Signature string
	Algorithm nullable.Nullable[string]
	Created   nullable.Nullable[string]
	Expires   nullable.Nullable[string]
	Headers   nullable.Nullable[string]
}

var _ fmt.Stringer = SignatureParams{}

func simpleQuotes(str string) string {
	return "\"" + str + "\""
}

func (sp SignatureParams) String() string {
	result := []string{}
	if sp.KeyID.HasValue() {
		result = append(result, KeyID+"="+simpleQuotes(sp.KeyID.ValueOrDefault("")))
	}
	result = append(result, Signature+"="+sp.Signature)
	if sp.Algorithm.HasValue() {
		result = append(result, Algorithm+"="+simpleQuotes(sp.Algorithm.ValueOrDefault("")))
	}
	if sp.Created.HasValue() {
		result = append(result, Created+"="+simpleQuotes(sp.Created.ValueOrDefault("")))
	}
	if sp.Expires.HasValue() {
		result = append(result, Expires+"="+simpleQuotes(sp.Expires.ValueOrDefault("")))
	}
	if sp.Headers.HasValue() {
		result = append(result, Headers+"="+simpleQuotes(sp.Headers.ValueOrDefault("")))
	}
	return strings.Join(result, ", ")
}

func ParseSignatureParams(params string) SignatureParams {
	result := SignatureParams{}
	for _, param := range strings.Split(params, ",") {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			continue
		}
		fieldName := strings.TrimSpace(parts[0])
		fieldValue := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		switch fieldName {
		case KeyID:
			result.KeyID = nullable.Just(fieldValue)
		case Signature:
			result.Signature = fieldValue
		case Algorithm:
			result.Algorithm = nullable.Just(fieldValue)
		case Created:
			result.Created = nullable.Just(fieldValue)
		case Expires:
			result.Expires = nullable.Just(fieldValue)
		case Headers:
			result.Headers = nullable.Just(fieldValue)
		}
	}
	return result
}
