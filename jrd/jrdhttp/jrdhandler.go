package jrdhttp

import (
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"net/http"
)

func CreateJRDHandler(handler func(r *http.Request) (jrd.JRD, httperrors.HTTPError)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j, errjrd := handler(r)
		if errjrd != nil {
			errjrd.ServeHTTP(w, r)
			return
		}
		err := httphelpers.WriteJSON(w, j)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		}
	})
}
