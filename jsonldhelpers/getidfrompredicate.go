package jsonldhelpers

import "fediverse/slices"

// GetIDFromPredicate gets a single ID, given the predicate to an object.
//
// If more than one object is associated with the predicate, then this function
// will return false.
//
// It is strongly recommended that the supplied object be of an expanded formm.
func GetIDFromPredicate(v any, predicate string) (string, bool) {
	objects := GetObjects(v, predicate)

	if len(objects) != 1 {
		return "", false
	}

	object, ok := slices.First(objects)
	if !ok {
		return "", false
	}

	msa, ok := object.(map[string]any)
	if !ok {
		return "", false
	}

	idAny, ok := msa["@id"]
	if !ok {
		return "", false
	}

	id, ok := idAny.(string)
	return id, ok
}
