package activitystreams

import (
	"encoding/json"
	"fediverse/maybe"
)

type Document struct {
	Object
}

var _ json.Marshaler = Link{}

func (d Document) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(d)
}
