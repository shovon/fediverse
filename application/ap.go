package application

import (
	"fediverse/application/config"
	"fediverse/application/lib"
	"fediverse/functional"
	"fediverse/httphelpers"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/nullable"
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

func ap() func(http.Handler) http.Handler {
	return hh.Accept([]string{"application/*+json"}).
		Process(hh.Group("/ap",
			hh.Group(
				"/users/:username",
				functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
					func(next http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							if lib.UserExists(httphelpers.GetRouteParam(r, "username")) {
								w.WriteHeader(404)
								w.Write([]byte("Not Found"))
								return
							}

							next.ServeHTTP(w, r)
						})
					},
					hh.Processors{
						hh.Method("GET"),
						hh.Route("/"),
					}.Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						a := func(path string) possibleerror.PossibleError[string] {
							u, err := requestbaseurl.GetRequestURL(r)
							if err != nil {
								return possibleerror.Error[string](err)
							}
							return resolveURIToString(u.ResolveReference(r.URL), path)
						}

						err := hh.WriteJSON(w, 200, map[string]interface{}{
							jsonldkeywords.Context: []interface{}{
								"https://www.w3.org/ns/activitystreams",
							},
							"id":                        baseURL().ResolveReference(r.URL).String(),
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
						}, nullable.Just("application/activty+json; charset=utf-8"))
						if err != nil {
							httperrors.InternalServerError().ServeHTTP(w, r)
							return
						}
					}))),
					hh.Processors{
						hh.Method("GET"),
						hh.Route("/followers"),
					}.Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(501)
						w.Write([]byte("Not implemented"))
					}))),
					hh.Processors{
						hh.Method("POST"),
						hh.Route("/inbox"),
					}.Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(501)
						w.Write([]byte("Not implemented"))
					}))),
				}),
			),
		))
}
