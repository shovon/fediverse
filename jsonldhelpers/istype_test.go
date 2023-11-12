package jsonldhelpers

import "testing"

const documentType = "https://example.com/schema/ns#Document"

func TestIsType(t *testing.T) {
	var apType any = map[string]any{"@type": []any{documentType}}

	if !IsType(apType, documentType) {
		t.Errorf("Expected %v to be determined to be of type %s", apType, documentType)
		t.Fail()
	}
}
