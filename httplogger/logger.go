package httplogger

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.bytesWritten += int64(len(b))
	return lrw.ResponseWriter.Write(b)
}

func formatBytes(bytes int) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

	if bytes < 1 {
		return "0B"
	}

	b := float64(bytes)

	// Calculate the appropriate unit to use (e.g., KB, MB, etc.)
	exp := int(math.Log(b) / math.Log(1024))

	// Format the number with the unit prefix
	formatted := b / math.Pow(1024, float64(exp))

	return fmt.Sprintf("%.2f%s", formatted, units[exp])
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		start := time.Now()
		h.ServeHTTP(lrw, r)

		protocol := "http"
		if r.TLS != nil {
			protocol = "https"
		}

		fmt.Printf(
			"%s \u001b[35m%s \u001b[36m%s://%s%s %s \033[39m from %s - \033[32m%d \u001b[38;5;14m%s in \033[39m%v\n",
			time.Now().Format("2006-10-11 15:04:05"),
			r.Method,
			protocol,
			r.Host,
			r.URL.String(),
			r.Proto,
			r.RemoteAddr,
			lrw.statusCode,
			formatBytes(int(lrw.bytesWritten)),
			time.Since(start),
		)
	})
}
