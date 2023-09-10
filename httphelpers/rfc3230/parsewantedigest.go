package rfc3230

import (
	"fediverse/pair"
	"fediverse/slices"
	"strings"

	"github.com/shopspring/decimal"
)

func ParseWantedDigest(wantedDigest string) ([]pair.Pair[string, decimal.Decimal], error) {
	list := slices.Map(strings.Split(wantedDigest, ","), func(str string, _ int) string {
		return strings.TrimSpace(str)
	})

	pairs := []pair.Pair[string, decimal.Decimal]{}
	for _, str := range list {
		parts := strings.Split(str, ";")
		if len(parts) == 1 {
			pairs = append(pairs, pair.Pair[string, decimal.Decimal]{Left: parts[0], Right: decimal.NewFromInt(1)})
			continue
		}
		if len(parts) != 2 {
			return nil, MalformedDigestError()
		}

		qParts := strings.Split(parts[1], "=")
		if len(qParts) != 2 {
			return nil, MalformedDigestError()
		}

		if qParts[0] != "q" {
			return nil, MalformedDigestError()
		}

		qValue, err := decimal.NewFromString(qParts[1])
		if err != nil {
			return nil, MalformedDigestError()
		}

		pairs = append(pairs, pair.Pair[string, decimal.Decimal]{Left: parts[0], Right: qValue})
	}

	return pairs, nil
}
