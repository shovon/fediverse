package httphelpers

import (
	"context"
	"fediverse/nullable"
	"net/http"

	contentnegotiation "gitlab.com/jamietanna/content-negotiation-go"
)

type negotiatedMediaType struct{}

func Accept(supported []string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			negotiator := contentnegotiation.NewNegotiator(supported...)
			acceptHeader := r.Header.Get("Accept")
			negotiated, provided, err := negotiator.Negotiate(acceptHeader)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			middleware(next).
				ServeHTTP(
					w,
					r.WithContext(
						context.WithValue(
							r.Context(),
							negotiatedMediaType{},
							NegotiatedMediaType{negotiated, provided},
						),
					),
				)
		})
	}
}

func NotAccept(supported []string, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			negotiator := contentnegotiation.NewNegotiator(supported...)
			acceptHeader := r.Header.Get("Accept")
			_, _, err := negotiator.Negotiate(acceptHeader)
			if err == nil {
				next.ServeHTTP(w, r)
				return
			}

			middleware(next).ServeHTTP(w, r)
		})
	}
}

type NegotiatedMediaType struct {
	Negotiated contentnegotiation.MediaType
	Provided   contentnegotiation.MediaType
}

func GetParsedAccept(r *http.Request) nullable.Nullable[NegotiatedMediaType] {
	if n, ok := r.Context().Value(negotiatedMediaType{}).(NegotiatedMediaType); ok {
		return nullable.Just(n)
	}
	return nullable.Null[NegotiatedMediaType]()
}
