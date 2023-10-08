package server

import (
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
	"fediverse/pathhelpers"
	"fediverse/security/rsahelpers"
	"net/http"
	"os"
)

type Following string
type Follower string

func searchUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lib.UserExists(hh.GetRouteParam(r, "username")) {
			httperrors.NotFound().ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func actor() func(http.Handler) http.Handler {
	return functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
		hh.PartialRoute(userRoute, searchUser),
		hh.Processors{
			hh.Method("GET"),
			hh.Route(userRoute),
		}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
			key := keymanager.GetPrivateKey()

			// TODO: this should ideally be cached.
			pubKeyString, err := rsahelpers.PublicKeyToPKIXString(&key.PublicKey)
			if err != nil {
				return nil, httperrors.InternalServerError()
			}

			origin := requestbaseurl.GetRequestOrigin(r)

			params := map[string]string{
				// TODO: this should be soft-coded.
				"username": hh.GetRouteParam(r, "username"),
			}

			actorRoot := origin + pathhelpers.FillFields(userRoute, params)

			return map[string]any{
				"@context": []interface{}{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},

				"id": actorRoot,

				"type": "Person",

				// TODO: this should be soft-coded. That is, retrieve the username,
				//   given some lookup invocation.
				"preferredUsername": config.Username(),

				// TODO: this should be soft-coded. That is, retrieve the display name,
				//   given some lookup invocation.
				"name": config.DisplayName(),

				// TODO: also find a way to soft code this.
				"summary": "<p>This person doesn't have a bio yet.</p>",

				"following": origin + pathhelpers.FillFields(followingRoute, params),
				"followers": origin + pathhelpers.FillFields(followersRoute, params),
				"inbox":     origin + pathhelpers.FillFields(inboxRoute, params),
				"outbox":    origin + pathhelpers.FillFields(outboxRoute, params),
				"liked":     origin + pathhelpers.FillFields(likedRoute, params),

				// TODO: manually approving followers is definitely an important
				//   feature.
				"manuallyApprovesFollowers": false,
				"publicKey": map[string]any{
					"id":           actorRoot + "#main-key",
					"owner":        actorRoot,
					"publicKeyPem": pubKeyString,
				},

				// TODO:
				"endpoints": map[string]any{
					"sharedInbox": origin + sharedInbox,
				},
			}, nil
		}))),
		orderedcollection.Middleware(
			followersRoute,
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
			followersRoute,
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
			hh.Route(inboxRoute),
		}.Process(printbody.Middleware(os.Stdout)),
	})
}
