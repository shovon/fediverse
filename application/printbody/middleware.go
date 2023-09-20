package printbody

import (
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Middleware(writer io.Writer) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := httphelpers.ProcessBody(r)
			if err != nil {
				httperrors.InternalServerError().ServeHTTP(w, r)
				return
			}
			str := string(body)
			fmt.Fprint(os.Stdout, str) // Yes, we are ignoring the error here
			handler.ServeHTTP(w, r)
		})
	}
}
