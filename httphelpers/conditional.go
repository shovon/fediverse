package httphelpers

import (
	"net/http"
	"net/url"
)

type ReadOnlyHeader struct {
	header http.Header
}

func NewReadOnlyHeader(header http.Header) ReadOnlyHeader {
	return ReadOnlyHeader{header: header}
}

func (r ReadOnlyHeader) Values(key string) []string {
	return r.header.Values(key)
}

type ReadOnlyRequest struct {
	request *http.Request
}

func NewReadOnlyRequest(r *http.Request) ReadOnlyRequest {
	return ReadOnlyRequest{request: r}
}

func (r ReadOnlyRequest) Method() string {
	return r.request.Method
}

func (r ReadOnlyRequest) URL() *url.URL {
	return r.request.URL
}

func (r ReadOnlyRequest) Header() ReadOnlyHeader {
	return NewReadOnlyHeader(r.request.Header)
}

func Condition(predicate func(ReadOnlyRequest) bool, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			readOnlyRequest := NewReadOnlyRequest(r)
			if predicate(readOnlyRequest) {
				middleware(next).ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
