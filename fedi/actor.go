package fedi

import (
	"encoding/json"
	"errors"
	"fediverse/acct"
	"fediverse/activitystreams/asvocab"
	"fediverse/webfinger"
	"io"
	"net/http"

	"github.com/piprate/json-gold/ld"
)

type Account struct {
	User string
	Host string
}

func FetchActorDocument(account Account) (map[string]any, error) {
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

			var output map[string]any
			err = json.Unmarshal(fullLD, &output)
			if err != nil {
				return nil, err
			}

			return output, nil
		}
	}

	return nil, errors.New("no self links found")
}

func GetInboxIRI(actor map[string]any) (string, error) {
	// General idea is this, given the actor document, that is hopefully in
	// JSON-LD form, we should first expand the document, and ensure that the
	// document has only one "entry" or "root" node. Then we

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	expanded, err := proc.Expand(actor, options)
	if err != nil {
		if inboxID, ok := actor[asvocab.Inbox]; ok {
			if str, ok := inboxID.(string); ok {
				return str, nil
			}
			return "", errors.New("inbox ID is not a string")
		} else {
			return "", errors.New("no inbox ID found")
		}
	}
	for _, obj := range expanded {

	}
	return "", errors.New("no inbox ID found")
}

func Send(recipient Account, message map[string]string) error {
	actor, err := FetchActorDocument(recipient)
	if err != nil {
		return err
	}

}
