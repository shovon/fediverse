package httphelpers

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(
	w http.ResponseWriter,
	content any,
) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}
