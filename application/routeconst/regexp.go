package routeconst

import (
	"fediverse/application/memoizer"
	"regexp"
)

var cache = memoizer.Cache[string, *regexp.Regexp]{}

func compileRegexp(r string) *regexp.Regexp {
	return cache.Memoize(r, func(r string) *regexp.Regexp {
		return regexp.MustCompile(r)
	})
}
