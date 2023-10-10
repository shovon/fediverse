package server

import (
	"fediverse/application/activity/server/orderedcollection"
	"fediverse/application/config"
	"fediverse/application/keymanager"
	"fediverse/application/lib"
	"fediverse/application/printbody"
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
	return hh.ApplyMiddlewares(hh.MiddlewaresList{
		// The main user route
		hh.Processors{
			hh.Method("GET"),
			hh.Route(UserRoute),
		}.Process(hh.ToMiddleware(searchUser(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
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

			actorRoot := origin + pathhelpers.FillFields(UserRoute, params)

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

				"following": origin + pathhelpers.FillFields(FollowingRoute, params),
				"followers": origin + pathhelpers.FillFields(FollowersRoute, params),
				"inbox":     origin + pathhelpers.FillFields(InboxRoute, params),
				"outbox":    origin + pathhelpers.FillFields(OutboxRoute, params),
				"liked":     origin + pathhelpers.FillFields(LikedRoute, params),

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
					"sharedInbox": origin + SharedInbox,
				},
			}, nil
		})))),

		// The followers collection.
		hh.Processors{
			hh.Route(FollowingRoute),
		}.Process(hh.ApplyMiddlewares(hh.MiddlewaresList{
			searchUser,
			orderedcollection.Middleware(
				orderedcollection.NewOrderedCollection[Following](
					func(hh.ReadOnlyRequest) uint64 {
						return 0
					},
					func(hh.ReadOnlyRequest, orderedcollection.ItemsFunctionParams) []Following {
						return []Following{}
					},
				),
			),
		})),

		// The following collection
		hh.Processors{
			hh.Route(FollowersRoute),
		}.Process(hh.ApplyMiddlewares(hh.MiddlewaresList{
			searchUser,
			orderedcollection.Middleware(
				orderedcollection.NewOrderedCollection[Following](
					func(hh.ReadOnlyRequest) uint64 {
						return 0
					},
					func(hh.ReadOnlyRequest, orderedcollection.ItemsFunctionParams) []Following {
						return []Following{}
					},
				),
			),
		})),

		// The inbox route.
		hh.Processors{
			hh.Route(InboxRoute),
		}.Process(printbody.Middleware(os.Stdout)),
	})
}
