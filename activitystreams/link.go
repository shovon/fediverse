package activitystreams

import (
	"encoding/json"
	"fediverse/maybe"
)

type Link struct {
	Href      maybe.Maybe[string]          `json:"href"`
	Rel       maybe.Maybe[string]          `json:"rel"`
	MediaType maybe.Maybe[string]          `json:"mediaType"`
	Name      maybe.Maybe[string]          `json:"name"`
	HrefLang  maybe.Maybe[string]          `json:"hreflang"`
	Height    maybe.Maybe[float64]         `json:"height"`
	Width     maybe.Maybe[float64]         `json:"width"`
	Preview   maybe.Maybe[json.RawMessage] `json:"preview"`
}

var _ json.Marshaler = Link{}

func (l Link) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(l)
}
