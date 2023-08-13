package urlhelpers

import (
	"encoding/json"
	"net/url"
)

type Marshalable url.URL

var _ json.Marshaler = Marshalable(url.URL{})

func (m Marshalable) MarshalJSON() ([]byte, error) {
	return json.Marshal((*url.URL(&m)).String())
}
