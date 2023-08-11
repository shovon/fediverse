package relaxediri

import (
	"encoding/json"
	"fediverse/iri"
)

// TODO: unit test this

type IRI string

var _ json.Marshaler = IRI("")

func (i IRI) MarshalJSON() ([]byte, error) {
	_, err := iri.ParseIRI(string(i))
	if err != nil {
		return nil, err
	}
	return json.Marshal(i)
}
