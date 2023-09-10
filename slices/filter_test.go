package slices

import "testing"

func TestFilter(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	result := Filter(l, func(i int, _ int) bool {
		return i%2 == 0
	})

	if len(result) != 2 {
		t.Errorf("expected 2 items, but got %d item(s) instead", len(result))
	}

	if result[0] != 2 {
		t.Errorf("expected 2, but got %d", result[0])
	}

	if result[1] != 4 {
		t.Errorf("expected 4, but got %d", result[1])
	}
}
