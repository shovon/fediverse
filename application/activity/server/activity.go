package server

import (
	"crypto/rsa"
	"fediverse/application/activity/routes"
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
	return hh.Processors{
		hh.Condition(hh.IsAcceptable([]string{"application/*+json"})),
		hh.DefaultResponseHeader("Content-Type", []string{"application/activity+json"}),
	}.Process(
		functional.RecursiveApply([](func(http.Handler) http.Handler){
			hh.Processors{
				hh.Method("POST"),
				hh.Route(routes.SharedInbox{}.Route().FullRoute()),
			}.Process(printbody.Middleware(os.Stdout)),
			hh.PartialRoute(routes.Actors{}.Actor().Route().ParameterizedRoute(), actor()),
		}),
	)
}
