package httphelpers

import (
	"context"
	"fediverse/pathhelpers"
	"net/http"
)

func RoutePrefix(route string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hasMatch, _, params := pathhelpers.PartialMatch(route, GetRelativePath(r))

			if !hasMatch {
				next.ServeHTTP(w, r)
				return
			}

			newR := r.WithContext(r.Context())
			for key, value := range params {
				newR = newR.WithContext(context.WithValue(newR.Context(), contextValue{key}, value))
			}

			middleware(next).ServeHTTP(w, newR)
		})
	}
}
