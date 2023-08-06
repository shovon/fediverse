package httphelpers

import (
	"encoding/json"
	"fediverse/nullable"
	"net/http"
)

func WriteJSON(
	w http.ResponseWriter,
	status int,
	content any,
	optionalContentType nullable.Nullable[string],
) error {
	w.
		Header().
		Add("Content-Type", optionalContentType.ValueOrDefault("application/json"))
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	w.Write(b)
	return nil
}
