package jrdhttp

import (
	"fediverse/httphelpers"
	"fediverse/jrd"
	"fediverse/jrd/jrdhttp/jrdhttperrors"
	"fediverse/nullable"
	"net/http"
)

func CreateJRDHandler(handler func(r *http.Request) (jrd.JRD, jrdhttperrors.JRDHttpError)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j, errjrd := handler(r)
		if errjrd != nil {
			errjrd.ServeHTTP(w, r)
			return
		}
		err := httphelpers.WriteJSON(w, 200, j, nullable.Just("application/jrd+json"))
		if err != nil {
			errjrd.ServeHTTP(w, r)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
	})
}
