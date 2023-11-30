package slices

import (
	"testing"
)

func TestHas(t *testing.T) {
	s := []string{"something"}

	if !Has(s, "something") {
		t.Log("Expected the set to claim it has 'something', but it does not")
		t.Logf("%v", s)
		t.Fail()
	}
}
