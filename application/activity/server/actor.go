package server

import (
	"fediverse/application/activity/routes"
	"fediverse/application/activity/server/orderedcollection"
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
	"fediverse/security/rsahelpers"
	"net/http"
	"os"
)

type Following string
type Follower string

func actor() func(http.Handler) http.Handler {
	return functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !lib.UserExists(hh.GetRouteParam(r, routes.ActorParam{}.ParameterName())) {
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

			key := keymanager.GetPrivateKey()
			pubKeyString, err := rsahelpers.PublicKeyToPKIXString(&key.PublicKey)
			if err != nil {
				return nil, httperrors.InternalServerError()
			}

			actorRoot := requestbaseurl.GetRequestOrigin(r) + r.URL.Path

			// Root: requestbaseurl.GetRequestOrigin(r)

			// id := a("")

			return map[string]any{
				jsonldkeywords.Context: []interface{}{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},
				"id":                        actorRoot,
				"type":                      "Person",
				"preferredUsername":         config.Username(),
				"name":                      config.DisplayName(),
				"summary":                   "This person doesn't have a bio yet.",
				"following":                 actorRoot + routes.Following{}.FullRoute(),
				"followers":                 actorRoot + "/followers",
				"inbox":                     actorRoot + "/inbox",
				"outbox":                    actorRoot + "/outbox",
				"liked":                     actorRoot + "/liked",
				"manuallyApprovesFollowers": false,
				"publicKey": map[string]any{
					"id":           actorRoot + "#main-key",
					"owner":        actorRoot,
					"publicKeyPem": pubKeyString,
				},
				"endpoints": map[string]any{
					"sharedInbox": actorRoot + "/" + routes.Root{}.Activity().SharedInbox().FullRoute(),
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
			hh.Route("/inbox"),
		}.Process(printbody.Middleware(os.Stdout)),
	})
}
