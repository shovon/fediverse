package httphelpers

import "net/http"

func SetHeader(key, value string) Processor {
	return ProcessorFunc(func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(key, value)
				middleware(next).ServeHTTP(w, r)
			})
		}
	})
}
