package jsonldcontext

import "encoding/json"

type ValidContext interface {
	uselessValidContext() useless
	MarshalJSON() ([]byte, error)
}

type testValidContext struct{}

var _ ValidContext = testValidContext{}
var _ json.Marshaler = testValidContext{}

func (t testValidContext) uselessValidContext() useless {
	return useless{}
}

func (t testValidContext) MarshalJSON() ([]byte, error) {
	return nil, nil
}
