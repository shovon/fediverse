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
	Count(hh.ReadOnlyRequest) uint64
	Items(hh.ReadOnlyRequest, ItemsFunctionParams) []V
}

type orderedCollectionRetriever[V any] struct {
	count func(hh.ReadOnlyRequest) uint64
	items func(hh.ReadOnlyRequest, ItemsFunctionParams) []V
}

var _ OrderedCollectionRetriever[any] = orderedCollectionRetriever[any]{}

func NewOrderedCollection[V any](
	count func(hh.ReadOnlyRequest) uint64,
	items func(hh.ReadOnlyRequest, ItemsFunctionParams) []V,
) OrderedCollectionRetriever[V] {
	return orderedCollectionRetriever[V]{count, items}
}

func (o orderedCollectionRetriever[V]) Count(req hh.ReadOnlyRequest) uint64 {
	return o.count(req)
}

func (o orderedCollectionRetriever[V]) Items(
	req hh.ReadOnlyRequest,
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
				hh.Condition(func(r hh.ReadOnlyRequest) bool {
					return strings.TrimSpace(r.URL.Query().Get("page")) != ""
				}),
				hh.ConditionMust(func(r hh.ReadOnlyRequest) bool {
					return isNaturalNumber(r.URL.Query().Get("page"))
				}, httperrors.BadRequest()),
			}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
				bbreq, err := hh.ToReadOnlyRequest(r)
				if err != nil {
					return nil, err
				}
				count := retriever.Count(bbreq)

				root := requestbaseurl.GetRequestOrigin(r) + r.URL.Path

				return map[string]any{
					jsonldkeywords.Context: []interface{}{
						"https://www.w3.org/ns/activitystreams",
					},
					"id":         root,
					"type":       "OrderedCollectionPage",
					"totalItems": count,
					"partOf":     root,

					// TODO: stuff here.
				}, nil
			}))),
			hh.Processors{
				hh.Method("GET"),
				hh.Route("/"),
			}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
				bbreq, err := hh.ToReadOnlyRequest(r)
				if err != nil {
					return nil, err
				}
				count := retriever.Count(bbreq)

				root := requestbaseurl.GetRequestOrigin(r) + r.URL.Path

				document := map[string]any{
					jsonldkeywords.Context: []interface{}{
						"https://www.w3.org/ns/activitystreams",
					},
					"id":         root,
					"type":       "OrderedCollection",
					"totalItems": count,
				}

				if count > 0 {
					document["first"] = possibleerror.Then(possibleerror.New(url.Parse(root)), possibleerror.MapToThen(func(s *url.URL) string {
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
