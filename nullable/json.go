package nullable

import (
	"encoding/json"
	"reflect"
	"strings"
)

func MarshalJSONWithMaybe(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)

	var sb strings.Builder
	sb.WriteRune('{')
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)
		tag := fieldType.Tag.Get("json")

		// Skip fields that are marked as "-"
		if tag == "-" {
			continue
		}

		fieldName := fieldType.Name

		// Use the field tag if available
		if tag != "" {
			tagValue := parseTag(tag)
			if tagValue != "" {
				fieldName = tagValue
			}
		}

		// Check if the field is of type NullableInt
		if strings.HasPrefix(fieldType.Type.Name(), "Nullable[") {
			isSet := field.FieldByName("hasValue").Bool()
			if !isSet {
				continue // Ignore if IsSet is false
			}
			if i > 0 {
				sb.WriteRune(',')
			}
			valueData, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			sb.WriteString("\"")
			sb.WriteString(fieldName)
			sb.WriteString("\":")
			sb.WriteString(string(valueData))
		} else {

			valueData, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			if i > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString("\"")
			sb.WriteString(fieldName)
			sb.WriteString("\":")
			sb.WriteString(string(valueData))
		}
	}
	sb.WriteRune('}')
	return []byte(sb.String()), nil
}

func parseTag(tag string) string {
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
