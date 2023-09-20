package urlhelpers

import (
	"fediverse/possibleerror"
	"net/url"
	p "path"
	"strings"
)

func JoinPath(u *url.URL, path string) possibleerror.PossibleError[*url.URL] {
	cloned, err := url.Parse(u.String())
	if err != nil {
		return possibleerror.Error[*url.URL](err)
	}
	pathCloned, err := url.Parse(path)
	if err != nil {
		return possibleerror.Error[*url.URL](err)
	}
	cloned.Path = p.Join(cloned.Path, pathCloned.Path)
	if strings.TrimSpace(cloned.Fragment) == "" {
		cloned.Fragment = pathCloned.Fragment
	}
	return possibleerror.NotError(cloned)
}
