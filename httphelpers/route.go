package httphelpers

import (
	"context"
	"net/http"
	"strings"
)

type contextValue struct {
	key string
}

// TODO: use the pathhelpers library

func Route(route string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathSplit := strings.Split(r.URL.Path, "/")
			routeSplit := strings.Split(route, "/")

			if len(pathSplit) != len(routeSplit) {
				next.ServeHTTP(w, r)
				return
			}

			newR := r.WithContext(r.Context())

			for i, pathPart := range pathSplit {
				routePart := routeSplit[i]
				if strings.HasPrefix(routePart, ":") {
					newR = r.WithContext(context.WithValue(r.Context(), contextValue{routePart[1:]}, pathPart))
					continue
				}
				if pathPart != routePart {
					next.ServeHTTP(w, r)
					return
				}
			}

			middleware(next).ServeHTTP(w, newR)
		})
	}
}

func GetRouteParam(r *http.Request, key string) string {
	value := r.Context().Value(contextValue{key})
	if value == nil {
		return ""
	}
	return value.(string)
}
