package jsonld

func GetObject(v any, predicate string) (any, bool) {
	if v == nil {
		return nil, false
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}

	o, ok := m[predicate]
	return o, ok
}
