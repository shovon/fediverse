package requestbaseurl

import (
	"context"
	"net/http"
	"net/url"
)

type overriden struct{}

func Override(u *url.URL) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), overriden{}, u)))
		})
	}
}
