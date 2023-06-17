package simpleld

import (
	"encoding/json"
	"fediverse/maybe"
)

type ContextTerm struct {
	ID        maybe.Maybe[string] `json:"@id"`
	Base      maybe.Maybe[string] `json:"@base"`
	Container maybe.Maybe[string] `json:"@container"`
	Context   maybe.Maybe[string] `json:"@context"`
	Direction maybe.Maybe[string] `json:"@direction"`
	Graph     maybe.Maybe[string] `json:"@graph"`
	Import    maybe.Maybe[string] `json:"@import"`
	Include   maybe.Maybe[string] `json:"@include"`
	Index     maybe.Maybe[string] `json:"@index"`
	JSON      maybe.Maybe[string] `json:"@json"`
	Language  maybe.Maybe[string] `json:"@language"`
	List      maybe.Maybe[string] `json:"@list"`
	Nest      maybe.Maybe[string] `json:"@nest"`
	None      maybe.Maybe[string] `json:"@none"`
	Prefix    maybe.Maybe[string] `json:"@prefix"`
	Propagate maybe.Maybe[string] `json:"@propagate"`
	Protected maybe.Maybe[string] `json:"@protected"`
	Reverse   maybe.Maybe[string] `json:"@reverse"`
	Set       maybe.Maybe[string] `json:"@set"`
	Type      maybe.Maybe[string] `json:"@type"`
	Value     maybe.Maybe[string] `json:"@value"`
	Version   maybe.Maybe[string] `json:"@version"`
	Vocab     maybe.Maybe[string] `json:"@vocab"`
}

var _ json.Marshaler = ContextTerm{}

func (c ContextTerm) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(c)
}
