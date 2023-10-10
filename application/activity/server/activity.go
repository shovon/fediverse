package server

import (
	"crypto/rsa"
	"fediverse/application/printbody"
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/possibleerror"
	"fediverse/security/rsahelpers"
	"fediverse/urlhelpers"
	"net/http"
	"net/url"
	"os"
)

func resolveURIToString(u *url.URL, path string) possibleerror.PossibleError[string] {
	return possibleerror.Then(
		urlhelpers.JoinPath(u, path), possibleerror.MapToThen(urlhelpers.ToString),
	)
}

// TODO: perhaps throw all of this into a third-party in-memory store.
var keyStore map[string]*rsa.PrivateKey = map[string]*rsa.PrivateKey{}

func getKey(id string) (*rsa.PrivateKey, error) {
	if keyStore[id] == nil {
		key, err := rsahelpers.GenerateRSPrivateKey(2048)
		if err != nil {
			return nil, err
		}
		keyStore[id] = key
	}
	return keyStore[id], nil
}

func getPublicKeyPEMString(id string) (string, error) {
	key, err := getKey(id)
	if err != nil {
		return "", err
	}
	return rsahelpers.PublicKeyToPKIXString(&key.PublicKey)
}

func ActivityPub() func(http.Handler) http.Handler {
	// TODO: be able to handle `Accept` activities, in response to a `Follow`
	//   activity.
	//
	// Here is an example payload:
	//
	// {
	//   "@context": "https://www.w3.org/ns/activitystreams",
	//   "id": "https://techhub.social/users/manlycoffee#accepts/follows/1129830",
	//   "type": "Accept",
	//   "actor": "https://techhub.social/users/manlycoffee",
	//   "object": {
	//     "id": "https://feditest.salrahman.com/activity/actors/johndoe/following/1",
	//     "type": "Follow",
	//     "actor": "https://feditest.salrahman.com/activity/actors/johndoe",
	//     "object": "https://techhub.social/users/manlycoffee"
	//   }
	// }

	return hh.Processors{
		hh.Condition(hh.IsAcceptable([]string{"application/*+json"})),
		hh.DefaultResponseHeader("Content-Type", []string{"application/activity+json"}),
	}.Process(
		functional.RecursiveApply([](func(http.Handler) http.Handler){

			hh.Processors{
				hh.Method("POST"),
				hh.Route(SharedInbox),
			}.Process(printbody.Middleware(os.Stdout)),

			actor(),
		}),
	)
}
