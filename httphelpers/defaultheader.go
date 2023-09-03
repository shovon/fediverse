package httphelpers

import "net/http"

func DefaultResponseHeader(header string, values []string) Processor {
	return ProcessorFunc((func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer middleware(next).ServeHTTP(w, r)
				if len(values) == 0 {
					return
				}
				first, rest := values[0], values[1:]
				w.Header().Set(header, first)
				for _, value := range rest {
					w.Header().Add(header, value)
				}
			})
		}
	}))
}
