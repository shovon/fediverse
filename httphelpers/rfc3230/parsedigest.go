package rfc3230

import (
	"errors"
	"fediverse/pair"
	"strings"
)

var malformedDigestError = errors.New("malformed digest")

func MalformedDigestError() error {
	return malformedDigestError
}

func Parse(header string) ([]pair.Pair[string, string], error) {
	digests := strings.Split(header, ",")
	pairs := []pair.Pair[string, string]{}
	for _, digest := range digests {
		arr := strings.Split(digest, "=")
		if len(arr) < 2 {
			return []pair.Pair[string, string]{}, MalformedDigestError()
		}
		pairs = append(pairs, pair.Pair[string, string]{
			Left:  arr[0],
			Right: strings.Join(arr[1:], "")},
		)
	}
	return pairs, nil
}