package rfc3230

import (
	"fediverse/pair"
	"testing"

	"github.com/shopspring/decimal"
)

func TestDeriveWantDigest(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		str := DeriveWantDigests([]pair.Pair[string, decimal.Decimal]{})
		if str != "" {
			t.Error("expected an empty string")
		}
	})
	t.Run("single 1", func(t *testing.T) {
		str := DeriveWantDigests([]pair.Pair[string, decimal.Decimal]{
			{Left: "sha-256", Right: decimal.NewFromInt(1)},
		})
		if str != "sha-256" {
			t.Errorf("expected a string equaling sha-256, but got %s", str)
		}
	})
}
