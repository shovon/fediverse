package nullable

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	type testStruct struct {
		Number Nilable[int]    `json:"number"`
		String Nilable[string] `json:"string"`
	}
	test := testStruct{
		Number: Just[int](5),
		String: Just[string]("Hello"),
	}
	data, err := MarshalJSONWithNilable(test)
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
		Number Nilable[int]    `json:"number"`
		String Nilable[string] `json:"string"`
	}
	test := testStruct{
		Number: Nil[int](),
		String: Nil[string](),
	}
	data, err := MarshalJSONWithNilable(test)
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
		Number Nilable[int]    `json:"number,omitempty"`
		String Nilable[string] `json:"string,omitempty"`
	}
	test := testStruct{
		Number: Nil[int](),
		String: Nil[string](),
	}
	data, err := MarshalJSONWithNilable(test)
	if err != nil {
		t.Error(err)
	}
	expected := `{}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
