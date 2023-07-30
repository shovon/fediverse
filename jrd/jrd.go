package jrd

import (
	"encoding/json"
	"fediverse/nullable"
	"time"
)

type JRD struct {
	Expires    nullable.Nullable[time.Time]                            `json:"expires,omitempty"`
	Subject    string                                                  `json:"subject"`
	Aliases    nullable.Nullable[[]string]                             `json:"aliases"`
	Properties nullable.Nullable[map[string]nullable.Nullable[string]] `json:"properties"`
	Links      nullable.Nullable[[]Link]                               `json:"links"`
}

var _ json.Marshaler = JRD{}

func (j JRD) MarshalJSON() ([]byte, error) {
	return nullable.MarshalJSONWithMaybe(j)
}
