package simpleld

import (
	"encoding/json"
	"errors"
	"fediverse/maybe"
)

const (
	urlType       = 0
	termsListType = 1
)

type Context struct {
	contextType int
	url         maybe.Maybe[string]
	termsList   maybe.Maybe[map[string]ContextTerm]
}

func NewURLContext(url string) Context {
	return Context{
		contextType: urlType,
		url:         maybe.Just(url),
	}
}

func NewTermsListContext(termsList map[string]ContextTerm) Context {
	return Context{
		contextType: termsListType,
		termsList:   maybe.Just(termsList),
	}
}

var _ json.Marshaler = Context{}

func (c Context) MarshalJSON() ([]byte, error) {
	switch c.contextType {
	case urlType:
		return maybe.MarshalJSONWithMaybe(c.url)
	case termsListType:
		return maybe.MarshalJSONWithMaybe(c.termsList)
	default:
		return nil, errors.New("Unknown context type")
	}
}
