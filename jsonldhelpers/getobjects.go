package jsonldhelpers

func GetObjects(v any, predicate string) []any {
	if v == nil {
		return nil
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}

	o, ok := m[predicate]
	if !ok {
		return nil
	}

	list, ok := o.([]any)
	if !ok {
		return nil
	}

	return list
}
