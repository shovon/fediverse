package ap

import (
	"fediverse/httphelpers"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/possibleerror"
	"fmt"
	"net/http"
	"net/url"
)

type OrderedCollectionMeta struct {
	TotalItems int
}

func OrderedCollection(route string, handler func(req *http.Request) OrderedCollectionMeta) func(http.Handler) http.Handler {
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
				return possibleerror.NotError(u.ResolveReference(r.URL).JoinPath(path))
			}

			a := func(path string) possibleerror.PossibleError[string] {
				u, err := requestbaseurl.GetRequestURL(r)
				if err != nil {
					return possibleerror.Error[string](err)
				}
				return resolveURIToString(u.ResolveReference(r.URL), path)
			}

			meta := handler(r)

			id := a("")

			return map[string]any{
				jsonldkeywords.Context: []interface{}{
					"https://www.w3.org/ns/activitystreams",
				},
				"id":         id,
				"type":       "OrderedCollection",
				"totalItems": meta.TotalItems,
				"first": possibleerror.Then(u(""), possibleerror.MapToThen(func(s *url.URL) string {
					fmt.Println(s)
					v := s.Query()
					v.Add("page", "1")
					s.RawQuery = v.Encode()
					return s.String()
				})),
			}, nil
		}))),
	)
}
