package server

import (
	"fediverse/application/config"
	"fediverse/application/keymanager"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/jsonhelpers/jsonhttp"
	"fediverse/pathhelpers"
	"fediverse/security/rsahelpers"
	"net/http"
)

func actorRoute() func(http.Handler) http.Handler {
	return hh.Processors{
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

		// Meh, don't bother with manual compaction.

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
	}))))
}
