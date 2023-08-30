package rfc3230

import (
	"fediverse/pair"
	"net/http"
)

func AddDigestsToHeaders(h http.Header, body []byte, digesters []Digester) error {
	digestParts := []pair.Pair[string, string]{}

	for _, digester := range digesters {
		digest, err := digester.Digest(body)
		if err != nil {
			return err
		}
		digestParts = append(
			digestParts,
			pair.Pair[string, string]{Key: digester.Token(), Value: digest},
		)
	}

	digest := CreateDigestString(digestParts)

	h.Add("Digest", digest)

	return nil
}
