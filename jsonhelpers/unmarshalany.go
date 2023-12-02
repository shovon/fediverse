package jsonhelpers

import "encoding/json"

func UnmarshalAny(source []byte) (any, error) {
	var objOutput any
	err := json.Unmarshal(source, &objOutput)
	if err != nil {
		return nil, err
	}

	return objOutput, nil
}
