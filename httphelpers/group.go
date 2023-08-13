package httphelpers

import (
	"context"
	"fediverse/pathhelpers"
	"net/http"
	"path"
)

type relativePath struct{}

func Group(route string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hasMatch, remainder, params := pathhelpers.PartialMatch(route, GetRelativePath(r))

			if !hasMatch {
				next.ServeHTTP(w, r)
				return
			}

			newR := r.WithContext(r.Context())
			for key, value := range params {
				newR = newR.WithContext(context.WithValue(newR.Context(), contextValue{key}, value))
			}

			newR = newR.WithContext(context.WithValue(newR.Context(), relativePath{}, remainder))
			middleware(next).ServeHTTP(w, newR)
		})
	}
}

func GetRelativePath(r *http.Request) string {
	value := r.Context().Value(relativePath{})
	str, ok := value.(string)
	if ok {
		return path.Join("/", str)
	}
	return r.URL.Path
}
