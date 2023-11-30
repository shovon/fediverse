package jsonldhelpers

import (
	"fediverse/slices"
	"testing"
)

func TestDecode(t *testing.T) {
	t.Run("Decode flat", func(t *testing.T) {
		type testObject struct {
			ID   string                `mapstructure:"@id"`
			Name []SingleValue[string] `mapstructure:"https://example.com/ns#name"`
		}

		var data []testObject
		err := Decode([]byte(`{
		"@context": {
			"ex": "https://example.com/ns#",
			"name": "ex:name"
		},
		"@id": "https://example.com/application/people/1",
		"name": "John Doe"
	}`), &data)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
			t.Fail()
		}
		first, ok := slices.First(data)
		if !ok {
			t.Errorf("Expected slice to have at least one element, but it does not")
			t.FailNow()
		}
		if slices.FirstOrDefault(data, testObject{}).ID != "https://example.com/application/people/1" {
			t.Errorf("Expected ID to be https://example.com/application/people/1, but got %v", slices.FirstOrDefault(data, testObject{}).ID)
			t.Fail()
		}
		if slices.FirstOrDefault(first.Name, SingleValue[string]{Value: "Hello"}).Value != "John Doe" {
			t.Errorf("Expected name to be John Doe, but got %v", slices.FirstOrDefault(first.Name, SingleValue[string]{Value: "Hello"}).Value)
			t.Fail()
		}
	})

	t.Run("Decode nested", func(t *testing.T) {
		type nestedObject struct {
			Type []string              `mapstructure:"@type"`
			ID   string                `mapstructure:"@id"`
			Name []SingleValue[string] `mapstructure:"https://example.com/ns#name"`
		}

		type testObject struct {
			Type []string              `mapstructure:"@type"`
			ID   string                `mapstructure:"@id"`
			Name []SingleValue[string] `mapstructure:"https://example.com/ns#name"`
			Dogs []nestedObject        `mapstructure:"https://example.com/ns#dogs"`
		}

		var data []testObject
		err := Decode([]byte(`{
		"@context": {
			"ex": "https://example.com/ns#",
			"name": "ex:name",
			"Person": "ex:Person",
			"Dog": "ex:Dog",
			"dogs": {
				"@id": "ex:dogs",
				"@type": "@id"
			}
		},
		"@type": "Person",
		"@id": "https://example.com/application/people/1",
		"name": "John Doe",
		"dogs": {
			"@type": "Dog",
			"@id": "https://example.com/application/dogs/1",
			"name": "Waffles"
		}
	}`), &data)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
			t.Fail()
		}
		first, ok := slices.First(data)
		if !ok {
			t.Errorf("Expected slice to have at least one element, but it does not")
			t.FailNow()
		}
		if !slices.Has(first.Type, "https://example.com/ns#Person") {
			t.Errorf("Expected type to contain https://example.com/ns#Person, but it does not")
			t.FailNow()
		}
		if slices.FirstOrDefault(data, testObject{}).ID != "https://example.com/application/people/1" {
			t.Errorf("Expected ID to be https://example.com/application/people/1, but got %v", slices.FirstOrDefault(data, testObject{}).ID)
			t.Fail()
		}
		if slices.FirstOrDefault(first.Name, SingleValue[string]{Value: "Hello"}).Value != "John Doe" {
			t.Errorf("Expected name to be John Doe, but got %v", slices.FirstOrDefault(first.Name, SingleValue[string]{Value: "Hello"}).Value)
			t.Fail()
		}

		dog, ok := slices.First(first.Dogs)
		if !ok {
			t.Errorf("Expected at least one dog in the slice of dogs, but got none")
			t.Fail()
		}

		if !slices.Has(dog.Type, "https://example.com/ns#Dog") {
			t.Errorf("Expected type to contain https://example.com/ns#Dog, but it does not")
			t.Fail()
		}
		if dog.ID != "https://example.com/application/dogs/1" {
			t.Errorf("Expected ID to be https://example.com/application/dogs/1, but got %v", dog.ID)
			t.Fail()
		}
		if slices.FirstOrDefault(dog.Name, SingleValue[string]{}).Value != "Waffles" {
			t.Errorf("Expected name to be Waffles, but got %v", slices.FirstOrDefault(first.Dogs, nestedObject{}).Name)
			t.Fail()
		}
	})
}
