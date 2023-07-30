package nullable

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	type testStruct struct {
		Number Nullable[int]    `json:"number"`
		String Nullable[string] `json:"string"`
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

func TestNull(t *testing.T) {
	type testStruct struct {
		Number Nullable[int]    `json:"number"`
		String Nullable[string] `json:"string"`
	}
	test := testStruct{
		Number: Null[int](),
		String: Null[string](),
	}
	data, err := MarshalJSONWithMaybe(test)
	if err != nil {
		t.Error(err)
	}
	expected1 := `{"number":null,"string":null}`
	expected2 := `{"string":null,"number":null}`
	if string(data) != expected1 && string(data) != expected2 {
		t.Errorf("Expected %s, got %s", expected1, string(data))
	}
}

func TestMissing(t *testing.T) {
	type testStruct struct {
		Number Nullable[int]    `json:"number,omitempty"`
		String Nullable[string] `json:"string,omitempty"`
	}
	test := testStruct{
		Number: Null[int](),
		String: Null[string](),
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
