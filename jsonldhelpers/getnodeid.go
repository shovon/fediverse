package jsonldhelpers

func GetNodeID(v any) (string, bool) {
	if v == nil {
		return "", false
	}

	m, ok := v.(map[string]any)
	if !ok {
		return "", false
	}

	id, ok := m["@id"]
	if !ok {
		return "", false
	}

	idStr, ok := id.(string)
	return idStr, ok
}
