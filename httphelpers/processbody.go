package httphelpers

import (
	"bytes"
	"io"
	"net/http"
)

// ProcessBody reads the body of an HTTP request and returns it as a byte slice.
// Also replaces the http.Request.Body with a new io.ReadCloser that reads from
// the buffered byte slice.
func ProcessBody(r *http.Request) ([]byte, error) {
	byteSlice, err := io.ReadAll(r.Body)
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(byteSlice))
	return byteSlice, err
}
