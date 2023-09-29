package requestbaseurl

import (
	"net/http"
)

func GetRequestOrigin(r *http.Request) string {
	if u, ok := r.Context().Value(overriden{}).(string); ok {
		return u
	}

	if r.TLS != nil {
		return "https://" + r.Host
	}

	return "http://" + r.Host
}
