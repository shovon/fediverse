package maybe

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	type testStruct struct {
		Number Maybe[int]    `json:"number"`
		String Maybe[string] `json:"string"`
	}
	test := testStruct{
		Number: Just[int](5),
		String: Just[string]("Hello"),
	}
	data, err := MarshalJSONWithMaybe(test)
	if err != nil {
		t.Error(err)
	}
	expected := `{"number":5,"string":"Hello"}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestMissing(t *testing.T) {
	type testStruct struct {
		Number Maybe[int]    `json:"number"`
		String Maybe[string] `json:"string"`
	}
	test := testStruct{
		Number: Nothing[int](),
		String: Nothing[string](),
	}
	data, err := MarshalJSONWithMaybe(test)
	if err != nil {
		t.Error(err)
	}
	expected := `{}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
