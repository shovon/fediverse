package ap

import (
	"fediverse/application/config"
	"fediverse/application/lib"
	"fediverse/functional"
	"fediverse/httphelpers"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/possibleerror"
	"fediverse/urlhelpers"
	"net/http"
	"net/url"
)

func resolveURIToString(u *url.URL, path string) possibleerror.PossibleError[string] {
	return possibleerror.Then(
		urlhelpers.JoinPath(u, path), possibleerror.MapToThen(urlhelpers.ToString),
	)
}

func ActivityPub() func(http.Handler) http.Handler {
	return hh.Processors{
		hh.Accept([]string{"application/*+json"}),
		hh.DefaultHeader("Content-Type", []string{"application/activity+json"}),
	}.Process(hh.Group("/ap",
		hh.Group(
			"/users/:username",
			functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if !lib.UserExists(httphelpers.GetRouteParam(r, "username")) {
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

					return map[string]any{
						jsonldkeywords.Context: []interface{}{
							"https://www.w3.org/ns/activitystreams",
							"https://w3id.org/security/v1",
						},
						"id":                        a(""),
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
							"publicKeyPem": "",
						},
					}, nil
				}))),
				hh.Processors{
					hh.Method("GET"),
					hh.Route("/following"),
				}.Process(hh.ToMiddleware(httperrors.NotImplemented())),
				hh.Processors{
					hh.Method("GET"),
					hh.Route("/followers"),
				}.Process(hh.ToMiddleware(httperrors.NotImplemented())),
				hh.Processors{
					hh.Method("POST"),
					hh.Route("/inbox"),
				}.Process(hh.ToMiddleware(httperrors.NotImplemented())),
			}),
		),
	))
}
