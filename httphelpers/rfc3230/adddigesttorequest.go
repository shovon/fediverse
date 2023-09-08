package rfc3230

import (
	"fediverse/pair"
	"fediverse/slices"
	"io"
	"net/http"
	"strings"
)

// AddDigestToRequest adds a Digest header to an HTTP request.
func AddDigestToRequest(req *http.Request, digesters []Digester) error {
	slice, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(strings.NewReader(string(slice)))
	digests := []pair.Pair[string, string]{}
	for _, digester := range digesters {
		result, err := digester.Digest(slice)
		if err != nil {
			return err
		}
		digests = append(digests, pair.Pair[string, string]{Left: digester.Token(), Right: result})
	}
	s := strings.Join(slices.Map(digests, func(digest pair.Pair[string, string], _ int) string {
		return digest.Left + "=" + digest.Right
	}), ",")
	req.Header.Add("Digest", s)
	return nil
}
