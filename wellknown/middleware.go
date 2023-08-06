package wellknown

import (
	"fediverse/httphelpers"
	"fmt"
	"net/http"
)

func WellKnown(path string, handler http.Handler) func(http.Handler) http.Handler {
	return httphelpers.Route(fmt.Sprintf("./well-known/%s", path), httphelpers.ToHandlerFunc(handler))
}
