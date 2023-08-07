package httphelpers

import (
	"net/http"
)

type Middlewares []func(http.Handler) http.Handler

var _ http.Handler = Middlewares{}

func (m *Middlewares) Use(middleware func(http.Handler) http.Handler) {
	*m = append(*m, middleware)
}

func (m Middlewares) ToMiddleware() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		if len(m) == 0 {
			return handler
		}
		if len(m) == 1 {
			return m[0](handler)
		}
		return m[0](append(Middlewares(m[1:]), ToMiddleware(handler)))
	}
}

func (m Middlewares) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(m) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}
	if len(m) == 1 {
		m[0](http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
		})).ServeHTTP(w, r)
		return
	}
	m[0](Middlewares(m[1:])).ServeHTTP(w, r)
}
