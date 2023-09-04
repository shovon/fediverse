package orderedcollection

import (
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
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

type ItemsFunctionParams struct {
	PageNumber uint64
	MaxItems   uint64
}

// OrderedCollectionRetriever represents a series of methods that will allow you
// to get the total number of items in a collection, a subset of the list of
// items in the collection, given a page number and a limit.
type OrderedCollectionRetriever[V any] interface {
	Count(hh.BarebonesRequest) uint64
	Items(hh.BarebonesRequest, ItemsFunctionParams) []V
}

type orderedCollectionRetriever[V any] struct {
	count func(hh.BarebonesRequest) uint64
	items func(hh.BarebonesRequest, ItemsFunctionParams) []V
}

var _ OrderedCollectionRetriever[any] = orderedCollectionRetriever[any]{}

func NewOrderedCollection[V any](
	count func(hh.BarebonesRequest) uint64,
	items func(hh.BarebonesRequest, ItemsFunctionParams) []V,
) OrderedCollectionRetriever[V] {
	return orderedCollectionRetriever[V]{count, items}
}

func (o orderedCollectionRetriever[V]) Count(req hh.BarebonesRequest) uint64 {
	return o.count(req)
}

func (o orderedCollectionRetriever[V]) Items(
	req hh.BarebonesRequest,
	params ItemsFunctionParams,
) []V {
	return o.items(req, params)
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

func Middleware[V any](route string, retriever OrderedCollectionRetriever[V]) func(http.Handler) http.Handler {
	return hh.Group(
		route,
		functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
			hh.Processors{
				hh.Method("GET"),
				hh.Route("/"),
				hh.Condition(func(r hh.BarebonesRequest) bool {
					return strings.TrimSpace(r.URL.Query().Get("page")) != ""
				}),
				hh.ConditionMust(func(r hh.BarebonesRequest) bool {
					return isNaturalNumber(r.URL.Query().Get("page"))
				}, httperrors.BadRequest()),
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

				bbreq, err := hh.CopyRequest(r)
				if err != nil {
					return nil, err
				}
				count := retriever.Count(bbreq)

				id := a("")

				return map[string]any{
					jsonldkeywords.Context: []interface{}{
						"https://www.w3.org/ns/activitystreams",
					},
					"id":         id,
					"type":       "OrderedCollectionPage",
					"totalItems": count,
					"partOf": possibleerror.Then(u(""), possibleerror.MapToThen(func(u *url.URL) string {
						u.RawQuery = ""
						return u.String()
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

				bbreq, err := hh.CopyRequest(r)
				if err != nil {
					return nil, err
				}
				count := retriever.Count(bbreq)

				id := a("")

				document := map[string]any{
					jsonldkeywords.Context: []interface{}{
						"https://www.w3.org/ns/activitystreams",
					},
					"id":         id,
					"type":       "OrderedCollection",
					"totalItems": count,
				}

				if count > 0 {
					document["first"] = possibleerror.Then(u(""), possibleerror.MapToThen(func(s *url.URL) string {
						fmt.Println(s)
						v := s.Query()
						v.Add("page", "1")
						s.RawQuery = v.Encode()
						return s.String()
					}))
				}

				return document, nil
			}))),
		}),
	)
}
