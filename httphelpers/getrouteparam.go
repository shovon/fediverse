package httphelpers

import "net/http"

func GetRouteParam(r *http.Request, key string) string {
	value := r.Context().Value(contextValue{key})
	if value == nil {
		return ""
	}
	return value.(string)
}
