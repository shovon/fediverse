package httphelpers

import "net/http"

type Processor interface {
	Process(func(http.Handler) http.Handler) func(http.Handler) http.Handler
}

// Usage:
//
// Processor.Process(middleware)

// ProcessFunc is a function that accepts a middleware and returns another
// middleware. The purpose of this is the same purpose that middlewares have.
type ProcessorFunc func(func(http.Handler) http.Handler) func(http.Handler) http.Handler

var _ Processor = ProcessorFunc(nil)

func (p ProcessorFunc) Process(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return p(middleware)
}
