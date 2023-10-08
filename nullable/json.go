package nullable

import (
	"encoding/json"
	"reflect"
	"strings"
)

func MarshalJSONWithNilable(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)

	var sb strings.Builder
	sb.WriteRune('{')
	rt := rv.Type()
	pairs := make([]string, 0, rv.NumField())
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		jsonOptions := strings.Split(jsonTag, ",")
		var jsonName string
		var shouldOmitEmpty bool
		if len(jsonOptions) > 0 {
			jsonName = jsonOptions[0]
		}
		if len(jsonOptions) > 1 && jsonOptions[1] == "omitempty" {
			shouldOmitEmpty = true
		}

		// Skip fields that are marked as "-"
		if jsonName == "-" {
			continue
		}

		fieldName := fieldType.Name
		fieldType.Tag.Lookup("omitempty")

		// Use the field tag if available
		if jsonName != "" {
			tagValue := parseJSONTagName(jsonName)
			if tagValue != "" {
				fieldName = tagValue
			}
		}

		// Check if the field is of type Nullable
		if strings.HasPrefix(fieldType.Type.Name(), "Nilable[") {
			isSet := field.FieldByName("hasValue").Bool()
			if !isSet {
				if shouldOmitEmpty {
					continue
				}
				str, err := json.Marshal(fieldName)
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, string(str)+":null")
			} else {
				valueData, err := json.Marshal(field.Interface())
				if err != nil {
					return nil, err
				}
				str, err := json.Marshal(fieldName)
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, string(str)+":"+string(valueData))
			}
		} else {
			valueData, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			str, err := json.Marshal(fieldName)
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, string(str)+":"+string(valueData))
		}
	}
	sb.WriteString(strings.Join(pairs, ","))
	sb.WriteRune('}')
	return []byte(sb.String()), nil
}

func parseJSONTagName(tag string) string {
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
