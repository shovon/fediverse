package httphelpers

import (
	"context"
	"fediverse/pathhelpers"
	"net/http"
)

func Group(route string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hasMatch, remainder, params := pathhelpers.PartialMatch(route, r.URL.Path)

			if !hasMatch {
				next.ServeHTTP(w, r)
				return
			}

			newR := r.WithContext(r.Context())
			for key, value := range params {
				newR = newR.WithContext(context.WithValue(newR.Context(), contextValue{key}, value))
			}
			oldPath := newR.URL.Path
			newR.URL.Path = remainder
			middleware(next).ServeHTTP(w, newR)
			newR.URL.Path = oldPath
		})
	}
}
