package set

import "testing"

func TestHas(t *testing.T) {
	s := New("something")

	if !s.Has("something") {
		t.Log("Expected the set to claim it has 'something', but it does not")
		t.Logf("%v", s)
		t.Fail()
	}
}
