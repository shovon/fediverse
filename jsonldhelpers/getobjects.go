package jsonldhelpers

func GetObjects(v any, predicate string) ([]any, bool) {
	if v == nil {
		return nil, false
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}

	o, ok := m[predicate]
	if !ok {
		return nil, false
	}

	list, ok := o.([]any)
	return list, ok
}
