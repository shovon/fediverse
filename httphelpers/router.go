package httphelpers

import (
	"net/http"
)

type Router Middlewares

func (r *Router) Use(path string, handler http.Handler) {
	m := Middlewares(*r)
	m.Use(Route(path, handler))
	*r = Router(m)
}

func (r Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m := Middlewares(r)
	m.ServeHTTP(res, req)
}
