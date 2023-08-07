package httphelpers

import (
	"net/http"
)

func recursiveApply[T any](fn []func(T) T, d T) T {
	if len(fn) == 0 {
		return d
	}
	return fn[0](recursiveApply(fn[1:], d))
}

func MiddlewaresList(middlewares []func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return recursiveApply[http.Handler](middlewares, handler)
	}
}
