package fedi

import (
	"errors"
	"fediverse/acct"
	"fediverse/webfinger"
	"io"
	"net/http"
)

type Account struct {
	User string
	Host string
}

func GetActor(account Account) (any, error) {
	// First, we issue a WebFinger request to the specific host, requesting
	// specifically for the application/activity+json content type.

	j, err := webfinger.Lookup(account.Host, acct.Acct(account).String(), []string{"self"})
	if err != nil {
		return nil, err
	}

	links, ok := j.Links.Value()
	if !ok {
		return nil, errors.New("no links found")
	}

	for _, link := range links {
		if link.Rel == "self" {
			req, err := http.NewRequest("GET", link.Href, nil)
			if err != nil {
				return nil, err
			}
			fullLD, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			return fullLD, nil
			break
		}
	}

	return nil, errors.New("no self links found")
}
