package httphelpers

import "net/http"

type Processor interface {
	Process(func(http.Handler) http.Handler) func(http.Handler) http.Handler
}

type ProcessorFunc func(func(http.Handler) http.Handler) func(http.Handler) http.Handler

var _ Processor = ProcessorFunc(nil)

func (p ProcessorFunc) Process(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return p(middleware)
}
