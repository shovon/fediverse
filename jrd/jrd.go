package jrd

import (
	"encoding/json"
	"fediverse/nullable"
	"time"
)

type JRD struct {
	Expires    nullable.Nilable[time.Time]                           `json:"expires,omitempty"`
	Subject    nullable.Nilable[string]                              `json:"subject,omitempty"`
	Aliases    nullable.Nilable[[]string]                            `json:"aliases,omitempty"`
	Properties nullable.Nilable[map[string]nullable.Nilable[string]] `json:"properties,omitempty"`
	Links      nullable.Nilable[[]Link]                              `json:"links,omitempty"`
}

var _ json.Marshaler = JRD{}

func (j JRD) MarshalJSON() ([]byte, error) {
	return nullable.MarshalJSONWithNilable(j)
}
