package wellknown

import (
	"fediverse/httphelpers"
	"net/http"
)

func WellKnown(path string, handler http.Handler) func(http.Handler) http.Handler {
	return httphelpers.PartialRoute(
		"/.well-known",
		httphelpers.PartialRoute("/"+path, httphelpers.ToMiddleware(handler)),
	)
}
