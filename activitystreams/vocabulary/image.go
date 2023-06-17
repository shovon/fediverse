package activitystreams

import (
	"encoding/json"
	"fediverse/maybe"
)

type Image struct {
	Document
}

var _ json.Marshaler = Link{}

func (i Image) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(i)
}
