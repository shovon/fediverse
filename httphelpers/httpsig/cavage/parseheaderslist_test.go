package cavage

import (
	"fmt"
	"testing"
)

func TestParseHeadersList(t *testing.T) {
	headers := "date (request-target) host digest"
	expected := []string{"date", "(request-target)", "host", "digest"}

	actual, err := ParseHeadersList(headers)
	fmt.Println(actual, len(actual))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("expected %d headers, got %d", len(expected), len(actual))
	}

	for i, header := range actual {
		if header != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], header)
		}
	}
}
