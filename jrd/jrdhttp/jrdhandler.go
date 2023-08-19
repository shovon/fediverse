package jrdhttp

import (
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/json/jsonhttp"
	"net/http"
)

func CreateJRDHandler(handler func(r *http.Request) (jrd.JRD, httperrors.HTTPError)) http.Handler {
	return jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
		return handler(r)
	})
}
