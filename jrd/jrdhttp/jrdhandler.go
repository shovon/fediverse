package jrdhttp

import (
	"encoding/json"
	"fediverse/jrd"
	"net/http"
)

func CreateJRDHandler(j jrd.JRD) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		js, err := json.Marshal(j)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/jrd+json")
		w.Write(js)
	})
}
