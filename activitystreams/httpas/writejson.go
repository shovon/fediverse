package httpas

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	w.Header().Set("Content-Type", "application/activity+json; charset=utf-8")
	return json.NewEncoder(w).Encode(v)
}
