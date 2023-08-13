package urlhelpers

import (
	"fediverse/possibleerror"
	"net/url"
)

func ResolvePath(baseURL *url.URL, path string) possibleerror.PossibleError[*url.URL] {
	ref, err := url.Parse(path)
	if err != nil {
		return possibleerror.Error[*url.URL](err)
	}
	return possibleerror.Value(baseURL.ResolveReference(ref))
}
