package schema2p0

import (
	"encoding/json"
)

type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type UsersStats struct {
	Total          uint `json:"total"`
	ActiveHalfyear uint `json:"activeHalfyear"`
	ActiveMonth    uint `json:"activeMonth"`
}

type Usage struct {
	Users         UsersStats `json:"users"`
	LocalPosts    uint       `json:"localPosts"`
	LocalComments uint       `json:"localComments"`
}

type Services struct {
	InboundP  []string `json:"inbound"`
	OutboundP []string `json:"outbound"`
}

type Schema struct {
	Software          Software `json:"software"`
	Usage             Usage    `json:"usage"`
	OpenRegistrations bool     `json:"openRegistrations"`
}

var _ json.Marshaler = Schema{}

func (s Schema) MarshalJSON() ([]byte, error) {
	type marshaler struct {
		Version   string   `json:"version"`
		Protocols []string `json:"protocols"`
		Metadata  struct {
			ChatEnabled bool `json:"chat_enabled"`
		} `json:"metadata"`
		Software          Software `json:"software"`
		Usage             Usage    `json:"usage"`
		OpenRegistrations bool     `json:"openRegistrations"`
		Services          Services `json:"services"`
	}

	services := Services{
		InboundP:  []string{},
		OutboundP: []string{},
	}

	m := marshaler{
		Version:   "2.0",
		Protocols: []string{"activitypub"},
		Metadata: struct {
			ChatEnabled bool `json:"chat_enabled"`
		}{ChatEnabled: false},
		Software:          s.Software,
		Usage:             s.Usage,
		OpenRegistrations: s.OpenRegistrations,
		Services:          services,
	}

	return json.Marshal(m)
}
