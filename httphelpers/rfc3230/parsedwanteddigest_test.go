package rfc3230

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestParseWantedDigest(t *testing.T) {
	{
		t.Run("empty", func(t *testing.T) {
			digests, err := ParseWantedDigest("")
			if err != nil {
				t.Error("unexpected error")
			}
			if len(digests) != 0 {
				t.Error("expected no digests")
			}
		})
		t.Run("single", func(t *testing.T) {
			digests, err := ParseWantedDigest("sha-256")
			if err != nil {
				t.Error("unexpected error")
			}

			if len(digests) != 1 {
				t.Error("expected single digest")
			}

			if digests[0].Left != "sha-256" {
				t.Errorf("expected digest %s, got %s", "sha-256", digests[0].Left)
			}

			if !digests[0].Right.Equal(decimal.NewFromInt(1)) {
				t.Errorf("expected q value %s, got %s", "1", digests[0].Right.String())
			}
		})

		t.Run("single with q", func(t *testing.T) {
			digests, err := ParseWantedDigest("sha-256;q=0.5")
			if err != nil {
				t.Error("unexpected error")
			}

			if len(digests) != 1 {
				t.Error("expected single digest")
			}

			if digests[0].Left != "sha-256" {
				t.Errorf("expected digest %s, got %s", "sha-256", digests[0].Left)
			}

			v, err := decimal.NewFromString("0.5")
			if err != nil {
				panic(err)
			}
			if !digests[0].Right.Equal(v) {
				t.Errorf("expected q value %s, got %s", "1", digests[0].Right.String())
			}
		})

		t.Run("multiple", func(t *testing.T) {
			digests, err := ParseWantedDigest("sha-256,sha-512")
			if err != nil {
				t.Error("unexpected error")
			}

			if len(digests) != 2 {
				t.Error("expected 2 digests")
			}

			if digests[0].Left != "sha-256" {
				t.Errorf("expected digest %s, got %s", "sha-256", digests[0].Left)
			}

			if digests[1].Left != "sha-512" {
				t.Errorf("expected digest %s, got %s", "sha-512", digests[0].Left)
			}

			if !digests[0].Right.Equal(decimal.NewFromInt(1)) {
				t.Errorf("expected q value %s, got %s", "1", digests[0].Right.String())
			}

			if !digests[1].Right.Equal(decimal.NewFromInt(1)) {
				t.Errorf("expected q value %s, got %s", "1", digests[1].Right.String())
			}
		})

		t.Run("multiple with q", func(t *testing.T) {
			digests, err := ParseWantedDigest("sha-256;q=0.5,sha-512;q=0.5")
			if err != nil {
				t.Error("unexpected error")
			}

			if len(digests) != 2 {
				t.Error("expected 2 digests")
			}

			if digests[0].Left != "sha-256" {
				t.Errorf("expected digest %s, got %s", "sha-256", digests[0].Left)
			}

			if digests[1].Left != "sha-512" {
				t.Errorf("expected digest %s, got %s", "sha-512", digests[0].Left)
			}

			digest1, err := decimal.NewFromString("0.5")
			if err != nil {
				panic(err)
			}
			if !digests[0].Right.Equal(digest1) {
				t.Errorf("expected q value %s, got %s", "0.5", digests[0].Right.String())
			}

			digest2, err := decimal.NewFromString("0.5")
			if err != nil {
				panic(err)
			}
			if !digests[1].Right.Equal(digest2) {
				t.Errorf("expected q value %s, got %s", "0.5", digests[1].Right.String())
			}
		})

		t.Run("multiple with q and spaces", func(t *testing.T) {
			digests, err := ParseWantedDigest("sha-256;q=0.5, sha-512;q=0.5")
			if err != nil {
				t.Error("unexpected error")
			}

			if len(digests) != 2 {
				t.Error("expected 2 digests")
			}

			if digests[0].Left != "sha-256" {
				t.Errorf("expected digest %s, got %s", "sha-256", digests[0].Left)
			}

			if digests[1].Left != "sha-512" {
				t.Errorf("expected digest %s, got %s", "sha-512", digests[0].Left)
			}

			digest1, err := decimal.NewFromString("0.5")
			if err != nil {
				panic(err)
			}
			if !digests[0].Right.Equal(digest1) {
				t.Errorf("expected q value %s, got %s", "0.5", digests[0].Right.String())
			}

			digest2, err := decimal.NewFromString("0.5")
			if err != nil {
				panic(err)
			}
			if !digests[1].Right.Equal(digest2) {
				t.Errorf("expected q value %s, got %s", "0.5", digests[1].Right.String())
			}
		})

		t.Run("multiple with q and spaces and extra semicolon", func(t *testing.T) {
			_, err := ParseWantedDigest("sha-256;q=0.5, sha-512;q=0.5;")
			if err == nil {
				t.Error("unexpected no error")
			}
		})
	}
}
