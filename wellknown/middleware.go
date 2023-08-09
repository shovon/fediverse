package wellknown

import (
	"fediverse/httphelpers"
	"net/http"
)

func WellKnown(path string, handler http.Handler) func(http.Handler) http.Handler {
	return httphelpers.Route("/.well-known/"+path, httphelpers.ToMiddleware(handler))
}
