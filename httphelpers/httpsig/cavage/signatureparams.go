package cavage

import (
	"errors"
	"fediverse/nullable"
	"fediverse/slices"
	"fmt"
	"strconv"
	"strings"
	"time"
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

type Params struct {
	KeyID     nullable.Nullable[string]
	Algorithm nullable.Nullable[string]
	Created   time.Time
	Expires   nullable.Nullable[time.Time]
	Headers   nullable.Nullable[[]string]
}

type ParamsWithSignature struct {
	Params
	Signature string
}

var _ fmt.Stringer = ParamsWithSignature{}

func simpleQuotes(str string) string {
	return "\"" + str + "\""
}

func (sp ParamsWithSignature) String() string {
	result := []string{}
	if sp.KeyID.HasValue() {
		result = append(result, KeyID+"="+simpleQuotes(sp.KeyID.ValueOrDefault("")))
	}
	result = append(result, Signature+"="+simpleQuotes(sp.Signature))
	if sp.Algorithm.HasValue() {
		result = append(result, Algorithm+"="+simpleQuotes(sp.Algorithm.ValueOrDefault("")))
	}
	result = append(result, Created+"="+strconv.FormatInt(sp.Created.Unix(), 10))
	if t, ok := sp.Expires.Value(); ok {
		result = append(result, Expires+"="+strconv.FormatInt(t.Unix(), 10))
	}
	if sp.Headers.HasValue() {
		result = append(result, Headers+"="+simpleQuotes(strings.Join(sp.Headers.ValueOrDefault([]string{}), " ")))
	}
	return strings.Join(result, ", ")
}

var errInvalidCreatedField = errors.New("the \"created\" field is not a valid Unix timestamp")
var errInvalidExpiresField = errors.New("the \"expires\" field is not a valid Unix timestamp")

func ErrInvalidCreatedField() error {
	return errInvalidCreatedField
}

func ErrInvalidExpiresField() error {
	return errInvalidExpiresField
}

// OK, so here's the deal:
//
//   - if either the `created` or the `expires` key is not a valid integer, then
//     the signature param is overall invalid
//   - all errors returned by the function is a client error. If it is a server
//     error, then this is bad, and must absolutely be investigated
func ParseSignatureParams(params string) (ParamsWithSignature, error) {
	result := ParamsWithSignature{}
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
			value, err := strconv.ParseInt(fieldValue, 10, 64)
			if err != nil {
				return ParamsWithSignature{}, errInvalidCreatedField
			}
			result.Created = time.Unix(value, 0)
		case Expires:
			value, err := strconv.ParseInt(fieldValue, 10, 0)
			if err != nil {
				return ParamsWithSignature{}, errInvalidExpiresField
			}
			result.Expires = nullable.Just(time.Unix(value, 0))
		case Headers:
			result.Headers = nullable.Just(slices.Map(strings.Split(fieldValue, " "), slices.IgnoreIndex(strings.TrimSpace)))
		}
	}
	return result, nil
}
