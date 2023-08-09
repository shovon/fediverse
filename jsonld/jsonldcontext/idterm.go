package jsonldcontext

import "encoding/json"

type IDTerm string

var _ Term = IDTerm("")

func (i IDTerm) uselessTerm() useless {
	return useless{}
}

func (i IDTerm) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"@id":   i,
		"@type": "@id",
	})
}
