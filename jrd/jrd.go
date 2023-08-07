package jrd

import (
	"encoding/json"
	"fediverse/nullable"
	"time"
)

type JRD struct {
	Expires    nullable.Nullable[time.Time]                            `json:"expires,omitempty"`
	Subject    nullable.Nullable[string]                               `json:"subject,omitempty"`
	Aliases    nullable.Nullable[[]string]                             `json:"aliases,omitempty"`
	Properties nullable.Nullable[map[string]nullable.Nullable[string]] `json:"properties,omitempty"`
	Links      nullable.Nullable[[]Link]                               `json:"links,omitempty"`
}

var _ json.Marshaler = JRD{}

func (j JRD) MarshalJSON() ([]byte, error) {
	return nullable.MarshalJSONWithMaybe(j)
}
