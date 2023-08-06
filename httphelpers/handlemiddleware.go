package httphelpers

import "net/http"

func HandleMiddleware(f func(http.Handler) http.Handler, handler http.Handler) http.Handler {
	return f(handler)
}
