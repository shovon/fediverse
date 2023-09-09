package rfc3230

import "testing"

func TestParseDigestWithEqualSign(t *testing.T) {
	{
		headerDigest := "sha-256=abc=def"
		digests, err := ParseDigest(headerDigest)
		if err != nil {
			t.Error(err)
		}
		if len(digests) != 1 {
			t.Error("Expected 1 digest")
		}
		if digests[0].Left != "sha-256" {
			t.Error("Expected sha-256")
		}
		if digests[0].Right != "abc=def" {
			t.Error("Expected abc=def")
		}
	}

	{
		headerDigest := "sha-256=abc==def"
		digests, err := ParseDigest(headerDigest)
		if err != nil {
			t.Error(err)
		}
		if len(digests) != 1 {
			t.Error("Expected 1 digest")
		}
		if digests[0].Left != "sha-256" {
			t.Error("Expected sha-256")
		}
		if digests[0].Right != "abc==def" {
			t.Error("Expected abc==def")
		}
	}
}
