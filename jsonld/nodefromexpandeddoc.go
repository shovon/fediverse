package jsonld

type LDObject struct {
	node any
}

type Node struct {
	node any
}

func (a Node) GetID() (string, bool) {
	if msa, ok := a.node.(map[string]any); ok {
		if id, ok := msa["@id"]; ok {
			return id.(string), true
		}
	}
	return "", false
}

func GetObjectFromNode(node any, predicate string) []LDObject {
	if msa, ok := node.(map[string]any); ok {
		if object, ok := msa[predicate], ok {
			return object
		}
	}
	return []LDObject{}
}
