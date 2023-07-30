package jrd

import (
	"encoding/json"
	"fediverse/nullable"
)

type Link struct {
	Rel        string                                                  `json:"rel"`
	Href       string                                                  `json:"href"`
	Type       nullable.Nullable[string]                               `json:"type"`
	Titles     nullable.Nullable[map[string]string]                    `json:"titles"`
	Properties nullable.Nullable[map[string]nullable.Nullable[string]] `json:"properties"`
}

var _ json.Marshaler = Link{}

func (j Link) MarshalJSON() ([]byte, error) {
	return nullable.MarshalJSONWithMaybe(j)
}
