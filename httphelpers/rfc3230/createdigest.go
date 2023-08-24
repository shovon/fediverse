package rfc3230

import (
	"fediverse/pair"
	"fediverse/slices"
	"strings"
)

func CreateDigestString(digests []pair.Pair[string, string]) string {
	return strings.Join(slices.Map(digests, func(p pair.Pair[string, string]) string {
		return p.Key + "=" + p.Value
	}), ",")
}
