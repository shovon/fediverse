package httphelpers

import (
	"fediverse/functional"
	"net/http"
)

type MiddlewaresList [](func(http.Handler) http.Handler)

func ApplyMiddlewares(middlewares [](func(http.Handler) http.Handler)) func(http.Handler) http.Handler {
	return functional.RecursiveApply[http.Handler](middlewares)
}
