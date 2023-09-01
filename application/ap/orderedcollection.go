package ap

import (
	"fediverse/httphelpers"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/possibleerror"
	"fediverse/urlhelpers"
	"net/http"
	"net/url"
)

func OrderedCollection(route string, handler func()) func(http.Handler) http.Handler {
	return httphelpers.Group(
		route,
		hh.Processors{
			hh.Method("GET"),
			hh.Route("/"),
		}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
			u := func(path string) possibleerror.PossibleError[*url.URL] {
				u, err := requestbaseurl.GetRequestURL(r)
				if err != nil {
					return possibleerror.Error[*url.URL](err)
				}
				return urlhelpers.JoinPath(u, path)
			}

			a := func(path string) possibleerror.PossibleError[string] {
				u, err := requestbaseurl.GetRequestURL(r)
				if err != nil {
					return possibleerror.Error[string](err)
				}
				return resolveURIToString(u.ResolveReference(r.URL), path)
			}

			id := a("")

			return map[string]any{
				jsonldkeywords.Context: []interface{}{
					"https://www.w3.org/ns/activitystreams",
				},
				"id":         id,
				"type":       "OrderedCollection",
				"totalItems": 0,
				"first": possibleerror.Then(u(""), possibleerror.MapToThen(func(s *url.URL) string {
					v := s.Query()
					v.Add("page", "1")
					s.RawFragment = v.Encode()
					return s.String()
				})),
			}, nil
		}))),
	)
}
