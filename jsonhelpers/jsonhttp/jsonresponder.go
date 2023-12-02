package jsonhttp

import (
	"encoding/json"
	"fediverse/httphelpers/httperrors"
	"fediverse/str"
	"net/http"
)

// JSONResponder is used to create a HTTP
func JSONResponder(handler func(r *http.Request) (any, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := handler(r)
		if err != nil {
			switch e := err.(type) {
			case httperrors.HTTPError:
				e.ServeHTTP(w, r)
			default:
				httperrors.InternalServerError().ServeHTTP(w, r)
			}
			return
		}
		b, err := json.Marshal(v)
		if err != nil {
			httperrors.InternalServerError().ServeHTTP(w, r)
			return
		}
		if str.IsEmpty(w.Header().Get("Content-Type")) {
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(b)
	})
}
