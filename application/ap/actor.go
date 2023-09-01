package ap

import (
	"fediverse/application/ap/orderedcollection"
	"fediverse/application/config"
	"fediverse/application/lib"
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/possibleerror"
	"net/http"
)

type OrderedCollectionMeta struct {
	TotalItems int
}

type OrderedCollectionRetriever interface {
	Meta()
	Page()
}

func actor() func(http.Handler) http.Handler {
	return functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !lib.UserExists(hh.GetRouteParam(r, "username")) {
					httperrors.NotFound().ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
		hh.Processors{
			hh.Method("GET"),
			hh.Route("/"),
		}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
			a := func(path string) possibleerror.PossibleError[string] {
				u, err := requestbaseurl.GetRequestURL(r)
				if err != nil {
					return possibleerror.Error[string](err)
				}
				return resolveURIToString(u.ResolveReference(r.URL), path)
			}

			username := hh.GetRouteParam(r, "username")
			key, err := getPublicKeyPEMString(username)
			if err != nil {
				return nil, err
			}

			id := a("")

			return map[string]any{
				jsonldkeywords.Context: []interface{}{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				"id":                        id,
				"type":                      "Person",
				"preferredUsername":         config.Username(),
				"name":                      config.DisplayName(),
				"summary":                   "This person doesn't have a bio yet.",
				"following":                 a("following"),
				"followers":                 a("followers"),
				"inbox":                     a("inbox"),
				"outbox":                    a("outbox"),
				"liked":                     a("liked"),
				"manuallyApprovesFollowers": false,
				"publicKey": map[string]any{
					"id":           a("#main-key"),
					"owner":        a(""),
					"publicKeyPem": key,
				},
			}, nil
		}))),
		orderedcollection.Middleware(
			"/following",
			func(req *http.Request) orderedcollection.OrderedCollection {
				return orderedcollection.OrderedCollection{
					TotalItems: 0,
				}
			},
		),
		hh.Processors{
			hh.Method("GET"),
			hh.Route("/followers"),
		}.Process(hh.ToMiddleware(httperrors.NotImplemented())),
		hh.Processors{
			hh.Method("POST"),
			hh.Route("/inbox"),
		}.Process(hh.ToMiddleware(httperrors.NotImplemented())),
	})
}
