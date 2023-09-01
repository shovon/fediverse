package orderedcollection

import (
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/possibleerror"
	"fediverse/urlhelpers"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type OrderedCollection struct {
	TotalItems int
}

type OrderedCollectionRetriever interface {
	Meta()
	Page()
}

func resolveURIToString(u *url.URL, path string) possibleerror.PossibleError[string] {
	return possibleerror.Then(
		urlhelpers.JoinPath(u, path), possibleerror.MapToThen(urlhelpers.ToString),
	)
}

func isNaturalNumber(str string) bool {
	str = strings.TrimSpace(str)

	if str == "" {
		return false
	}

	for i, ch := range str {
		if i == 0 && ('1' <= ch && ch <= '9') {
			continue
		}
		if '0' <= ch && ch <= '9' {
			continue
		}
		return false
	}

	return true
}

func Middleware(route string, handler func(req *http.Request) OrderedCollection) func(http.Handler) http.Handler {
	return hh.Group(
		route,
		functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
			hh.Processors{
				hh.Method("GET"),
				hh.Route("/"),
				hh.Condition(func(r hh.BarebonesRequest) bool {
					return isNaturalNumber((r.URL.Query().Get("page")))
				}),
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
					"type":       "OrderedCollectionPage",
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
		}),
	)
}
