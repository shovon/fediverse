package pathhelpers

import "testing"

func TestMatch(t *testing.T) {
	t.Run("Match empty", func(t *testing.T) {
		result, msa := Match("/a/b/c", "/a/b/c")
		if !result {
			t.Fail()
		}
		if len(msa) != 0 {
			t.Fail()
		}
	})
	t.Run("Match with params", func(t *testing.T) {
		result, msa := Match("/a/:b/c", "/a/1/c")
		if !result {
			t.Fail()
		}
		if len(msa) != 1 {
			t.Fail()
		}
		if msa["b"] != "1" {
			t.Fail()
		}
	})
}
