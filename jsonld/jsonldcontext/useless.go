package jsonldcontext

import "encoding/json"

type useless struct{}

type uselessTerm struct{}

var _ Term = uselessTerm{}
var _ json.Marshaler = uselessTerm{}

func (u uselessTerm) uselessTerm() useless {
	return useless{}
}

func (u uselessTerm) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}
