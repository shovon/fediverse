package httphelpers

import (
	"net/http"
)

func Condition(predicate func(r ReadOnlyRequest) bool) Processor {
	return ProcessorFunc(func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqCopy, err := ToReadOnlyRequest(r)
				if err != nil {
					middleware(next).ServeHTTP(w, r)
					return
				}
				if predicate(reqCopy) {
					middleware(next).ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r)
			})
		}
	})
}
