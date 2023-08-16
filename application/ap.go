package application

import (
	"fediverse/application/config"
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/nullable"
	"fediverse/possibleerror"
	"fediverse/urlhelpers"
	"net/http"
	"net/url"
)

func ap() func(http.Handler) http.Handler {
	resolveURIToString := func(u *url.URL, path string) possibleerror.PossibleError[string] {
		return possibleerror.Then(
			urlhelpers.JoinPath(u, path), possibleerror.MapToThen(urlhelpers.ToString),
		)
	}

	return hh.Accept([]string{"application/*+json"}).
		Process(hh.Group("/ap",
			hh.Group(
				"/users/:username",
				functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler){
					hh.Processors{
						hh.Method("GET"),
						hh.Route("/"),
					}.Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						err := hh.WriteJSON(w, 200, map[string]interface{}{
							jsonldkeywords.Context: []interface{}{
								"https://www.w3.org/ns/activitystreams",
							},
							"id":                        baseURL().ResolveReference(r.URL).String(),
							"type":                      "Person",
							"preferredUsername":         config.Username(),
							"name":                      config.DisplayName(),
							"summary":                   "This person doesn't have a bio yet.",
							"following":                 resolveURIToString(baseURL().ResolveReference(r.URL), "following"),
							"followers":                 resolveURIToString(baseURL().ResolveReference(r.URL), "followers"),
							"inbox":                     resolveURIToString(baseURL().ResolveReference(r.URL), "inbox"),
							"outbox":                    resolveURIToString(baseURL().ResolveReference(r.URL), "outbox"),
							"liked":                     resolveURIToString(baseURL().ResolveReference(r.URL), "liked"),
							"manuallyApprovesFollowers": false,
						}, nullable.Just("application/activty+json; charset=utf-8"))
						if err != nil {
							w.WriteHeader(500)
							w.Write([]byte("Internal Server Error"))
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
