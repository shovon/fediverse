package activitystreams

import (
	"encoding/json"
	"fediverse/maybe"
)

type Object struct {
	Attachment   maybe.Maybe[string] `json:"attachment"`
	AttributedTo maybe.Maybe[string] `json:"attributedTo"`
	Audience     maybe.Maybe[string] `json:"audience"`
	Content      maybe.Maybe[string] `json:"content"`
	Context      maybe.Maybe[string] `json:"context"`
	Name         maybe.Maybe[string] `json:"name"`
	EndTime      maybe.Maybe[string] `json:"endTime"`
	Generator    maybe.Maybe[string] `json:"generator"`
	Icon         maybe.Maybe[string] `json:"icon"`
	Image        maybe.Maybe[string] `json:"image"`
}

var _ json.Marshaler = Link{}

func (o Object) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(o)
}
