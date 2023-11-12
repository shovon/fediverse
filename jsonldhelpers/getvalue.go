package jsonldhelpers

func GetValue[T any](v any) (T, bool) {
	var zero T
	if v == nil {
		return zero, false
	}

	m, ok := v.(map[string]any)
	if !ok {
		return zero, false
	}

	id, ok := m["@id"]
	if !ok {
		return zero, false
	}

	idStr, ok := id.(T)
	return idStr, ok
}
