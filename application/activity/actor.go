package activity

import (
	"fediverse/application/activity/orderedcollection"
	"fediverse/application/config"
	"fediverse/application/keymanager"
	"fediverse/application/lib"
	"fediverse/application/printbody"
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/possibleerror"
	"fediverse/security/rsahelpers"
	"fediverse/urlhelpers"
	"net/http"
	"net/url"
	"os"
)

type Following string
type Follower string

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

			// username := hh.GetRouteParam(r, "username")

			key := keymanager.GetPrivateKey()
			pubKeyString, err := rsahelpers.PublicKeyToPKIXString(&key.PublicKey)
			if err != nil {
				return nil, httperrors.InternalServerError()
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
					"publicKeyPem": pubKeyString,
				},
				"endpoints": map[string]any{
					"sharedInbox": possibleerror.Then(possibleerror.New(requestbaseurl.GetRequestURL(r)), func(u *url.URL) possibleerror.PossibleError[string] {
						return possibleerror.NotError(urlhelpers.ToString(u.ResolveReference(&url.URL{
							Path: "/sharedinbox",
						})))
					}),
				},
			}, nil
		}))),
		orderedcollection.Middleware(
			"/following",
			orderedcollection.NewOrderedCollection[Following](
				func(hh.ReadOnlyRequest) uint64 {
					return 0
				},
				func(hh.ReadOnlyRequest, orderedcollection.ItemsFunctionParams) []Following {
					return []Following{}
				},
			),
		),
		orderedcollection.Middleware(
			"/followers",
			orderedcollection.NewOrderedCollection[Follower](
				func(hh.ReadOnlyRequest) uint64 {
					return 0
				},
				func(hh.ReadOnlyRequest, orderedcollection.ItemsFunctionParams) []Follower {
					return []Follower{}
				},
			),
		),
		hh.Processors{
			hh.Method("POST"),
			hh.Route("/inbox"),
		}.Process(printbody.Middleware(os.Stdout)),
	})
}
