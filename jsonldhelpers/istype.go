package jsonld

import (
	"fediverse/slices"
)

// TODO: perhaps also check if the type is a slice of strings, and not just a
// slice of any?

// IsType checks if the given JSON-LD document is of the given type.
//
// Note: it is **strongly** recommended that the value that you supply as the
// first parameter v has been expanded into a map[string]any.
func IsType(v any, expectedType string) bool {
	if v == nil {
		return false
	}

	m, ok := v.(map[string]any)
	if !ok {
		return false
	}

	t, ok := m["@type"]
	if !ok {
		return false
	}

	s, ok := t.(string)
	if ok {
		return s == expectedType
	}

	list, ok := t.([]any)
	if !ok {
		return false
	}

	return slices.Some(list, func(s any) bool {
		if s, ok := s.(string); ok {
			return s == expectedType
		}
		return false
	})
}
