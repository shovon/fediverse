package printbody

import (
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fmt"
	"io"
	"net/http"
)

func Middleware(writer io.Writer) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := httphelpers.ProcessBody(r)
			str := string(body)
			fmt.Fprint(writer, str)
			if err != nil {
				httperrors.InternalServerError().ServeHTTP(w, r)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
