package httphelpers

import "net/http"

type Processors []Processor

var _ Processor = Processors(nil)

func (p Processors) Process(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	for _, processor := range p {
		middleware = processor.Process(middleware)
	}
	return middleware
}
