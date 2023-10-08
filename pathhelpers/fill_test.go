package pathhelpers

import "testing"

func TestFillFields(t *testing.T) {
	result := FillFields("/:cool", map[string]string{
		"cool": "beans",
	})

	if result != "/beans" {
		t.Errorf("Expected '/beans', got '%s'", result)
	}
}
