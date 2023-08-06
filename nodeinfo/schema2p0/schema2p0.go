package schema2p0

import (
	"encoding/json"
)

type Services struct {
	// if your platform is going to support some actual inbound and outbound, then
	// add them here.
}

func (s Services) MarshalJSON() ([]byte, error) {
	type marshaler struct {
		Inbound  []string `json:"inbound"`
		Outbound []string `json:"outbound"`
	}

	m := marshaler{}

	return json.Marshal(m)
}

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

type Schema struct {
	Software Software `json:"software"`
	Usage    Usage    `json:"users"`
}

var _ json.Marshaler = Schema{}

func (s Schema) MarshalJSON() ([]byte, error) {
	type marshaler struct {
		Version   string   `json:"version"`
		Protocols []string `json:"protocols"`
		Metadata  struct {
			ChatEnabled bool `json:"chat_enabled"`
		} `json:"metadata"`
	}

	m := marshaler{
		Version:   "2.0",
		Protocols: []string{"activitypub"},
	}

	return json.Marshal(m)
}
