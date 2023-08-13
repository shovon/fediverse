package urlhelpers

import (
	"fediverse/possibleerror"
	"net/url"
	p "path"
)

func JoinPath(u *url.URL, path string) possibleerror.PossibleError[*url.URL] {
	cloned, err := url.Parse(u.String())
	if err != nil {
		return possibleerror.Error[*url.URL](err)
	}
	cloned.Path = p.Join(cloned.Path, path)
	return possibleerror.Value(cloned)
}
