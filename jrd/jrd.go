package jrd

import (
	"encoding/json"
	"fediverse/maybe"
	"fediverse/nullable"
	"time"
)

type JRD struct {
	Expires    maybe.Maybe[time.Time]                            `json:"expires"`
	Subject    string                                            `json:"subject"`
	Aliases    maybe.Maybe[[]string]                             `json:"aliases"`
	Properties maybe.Maybe[map[string]nullable.Nullable[string]] `json:"properties"`
	Links      maybe.Maybe[[]Link]                               `json:"links"`
}

var _ json.Marshaler = JRD{}

func (j JRD) MarshalJSON() ([]byte, error) {
	return maybe.MarshalJSONWithMaybe(j)
}
