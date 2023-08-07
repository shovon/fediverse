package httphelpers

import (
	"context"
	"net/http"
	"strings"
)

func Group(route string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathSplit := strings.Split(r.URL.Path, "/")
			routeSplit := strings.Split(route, "/")

			if len(pathSplit) < len(routeSplit) {
				next.ServeHTTP(w, r)
				return
			}

			newR := r.WithContext(r.Context())
			newPathSplit := pathSplit[len(pathSplit):]

			for i, routePart := range routeSplit {
				pathPart := pathSplit[i]
				if strings.HasPrefix(routePart, ":") {
					newR = r.WithContext(context.WithValue(r.Context(), contextValue{routePart[1:]}, pathPart))
					continue
				}
				if pathPart != routePart {
					next.ServeHTTP(w, r)
					return
				}
				newPathSplit = newPathSplit[1:]
			}

			newR.URL.Path = "/" + strings.Join(newPathSplit, "/")

			middleware(next).ServeHTTP(w, newR)
		})
	}
}
