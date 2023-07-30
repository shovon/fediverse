package jrd

import (
	"encoding/json"
	"fediverse/maybe"
	"fediverse/nullable"
)

type Link struct {
	Rel        string                                            `json:"rel"`
	Href       string                                            `json:"href"`
	Type       maybe.Maybe[string]                               `json:"type"`
	Titles     maybe.Maybe[map[string]string]                    `json:"titles"`
	Properties maybe.Maybe[map[string]nullable.Nullable[string]] `json:"properties"`
}

var _ json.Marshaler = Link{}

func (j Link) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(j)
}
