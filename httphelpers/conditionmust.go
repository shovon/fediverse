package httphelpers

import (
	"fediverse/httphelpers/httperrors"
	"net/http"
)

func ConditionMust(predicate func(ReadOnlyRequest) bool, defaultHandler http.Handler) Processor {
	return ProcessorFunc(func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqCopy, err := ToReadOnlyRequest(r)
				if err != nil {
					httperrors.InternalServerError().ServeHTTP(w, r)
					return
				}
				if predicate(reqCopy) {
					middleware(next).ServeHTTP(w, r)
					return
				}
				defaultHandler.ServeHTTP(w, r)
			})
		}
	})
}
