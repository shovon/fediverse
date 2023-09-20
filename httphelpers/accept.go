package httphelpers

import (
	"fmt"
	"strings"

	contentnegotiation "gitlab.com/jamietanna/content-negotiation-go"
)

type negotiatedMediaType struct{}

func IsAcceptable(supported []string) func(r ReadOnlyRequest) bool {
	return func(r ReadOnlyRequest) bool {
		negotiator := contentnegotiation.NewNegotiator(supported...)
		acceptHeader := r.Header.Get("Accept")
		if strings.TrimSpace(acceptHeader) == "" {
			return true
		}
		_, _, err := negotiator.Negotiate(acceptHeader)
		fmt.Println(acceptHeader, err)
		return err == nil
	}
}
