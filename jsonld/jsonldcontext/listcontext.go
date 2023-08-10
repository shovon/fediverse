package jsonldcontext

import "encoding/json"

type ListContext []ValidContext

var _ ValidContext = ListContext{}

func (l ListContext) uselessValidContext() useless {
	return useless{}
}

func (l ListContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(l)
}
