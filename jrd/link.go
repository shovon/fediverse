package jrd

import (
	"encoding/json"
	"fediverse/nullable"
)

type Link struct {
	Rel        string                                                `json:"rel"`
	Href       string                                                `json:"href"`
	Type       nullable.Nilable[string]                              `json:"type,omitempty"`
	Titles     nullable.Nilable[map[string]string]                   `json:"titles,omitempty"`
	Properties nullable.Nilable[map[string]nullable.Nilable[string]] `json:"properties,omitempty"`
}

var _ json.Marshaler = Link{}

func (j Link) MarshalJSON() ([]byte, error) {
	return nullable.MarshalJSONWithNilable(j)
}
