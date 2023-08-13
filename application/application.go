package application

import (
	"fediverse/acct"
	"fediverse/application/config"
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/jsonld/jsonldkeywords"
	"fediverse/nodeinfo"
	"fediverse/nullable"
	"fediverse/pathhelpers"
	"fediverse/webfinger"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5/middleware"
)

func origin() string {
	return config.HttpProtocol() + "://" + config.Hostname()
}

type UserHost struct {
	Username string
	Host     string
}

const usernameParamKey = "username"

func parseURLResource(resource string) (UserHost, *url.URL, bool) {
	urlQuery, errURL := url.Parse(resource)

	if errURL != nil {
		return UserHost{}, urlQuery, false
	}

	handlers := pathhelpers.Handlers{}
	handlers["/users/:"+usernameParamKey] = func(p map[string]string) (bool, map[string]string) { return true, p }
	handlers["/:"+usernameParamKey] = func(params map[string]string) (bool, map[string]string) {
		if params[usernameParamKey][0] != '@' {
			return false, nil
		}
		return true, map[string]string{usernameParamKey: params[usernameParamKey][1:]}
	}

	match, params := handlers.Handle(urlQuery.Path)
	if !match {
		return UserHost{}, urlQuery, false
	}

	return UserHost{params[usernameParamKey], urlQuery.Host}, urlQuery, true
}

func Start() {
	m := [](func(http.Handler) http.Handler){}

	// TODO: move away from Chi, and use some other logger library.
	m = append(m, middleware.Logger)

	m = append(m, webfinger.WebFinger(func(resource string) (jrd.JRD, httperrors.HTTPError) {
		acct, acctErr := acct.ParseAcct(resource)
		userHost, urlQuery, urlIsValid := parseURLResource(resource)

		if acctErr != nil && !urlIsValid {
			return jrd.JRD{}, httperrors.BadRequest()
		}

		var user, host string

		if acctErr == nil {
			user = acct.User
			host = acct.Host
		} else if urlIsValid {
			user = userHost.Username
			host = userHost.Host
		}

		// Note: if this software is a multi-tenant instance, then we will need to
		// check if the host exists in the database, and that it matches the HTTP
		// HOST header.
		if host != config.Hostname() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		// NOTE: if this software is a multi-user instance, then we will need to
		// check if the user exists in the database.
		//
		// If the software is also a multi-tenant instance, then we will need to
		// check that not only is the user in the database, but if it is associated
		// with the correct host.
		if user != config.Username() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		if urlIsValid && urlQuery.Scheme != config.HttpProtocol() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		return webFingerJRD(UserHost{user, host}), nil
	}))

	m = append(m, nodeinfo.CreateNodeInfoMiddleware(origin(), "/nodinfo", func() nodeinfo.NodeInfoProps {
		return nodeinfo.NodeInfoProps{
			Software: nodeinfo.SoftwareInfo{
				Name:    "fediverse",
				Version: "0.0.1",
			},
			OpenRegistrations: false,
			Usage: nodeinfo.Usage{
				Users: nodeinfo.UsersStats{
					Total: 1,

					// TODO: actually get all activie users in the last 6 months
					ActiveHalfyear: 0,

					// TODO: actually get all activie users in the last 30 days
					ActiveMonth: 0,
				},
				LocalPosts:    0, // TODO: actually get all local posts
				LocalComments: 0, // TODO: actually get all local comments
			},
		}
	}))

	activitypubMediaTypes := []string{"application/json", "application/activity+json"}

	m = append(m, hh.NotAccept(activitypubMediaTypes).Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Just an article. Coming soon"))
	}))))

	m = append(
		m,
		hh.Group("/ap",
			hh.Processors{
				hh.Accept(activitypubMediaTypes),
				hh.Method("GET"),
				hh.Route("/users/:username"),
			}.Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// TODO; log the error output from WriteJSON
				hh.WriteJSON(w, 200, map[string]interface{}{
					jsonldkeywords.Context: []interface{}{
						"https://www.w3.org/ns/activitystreams",
					},
					"id":                        origin() + "/users/" + config.Username(),
					"type":                      "Person",
					"preferredUsername":         config.Username(),
					"name":                      config.DisplayName(),
					"following":                 origin() + "/users/" + config.Username() + "/following",
					"followers":                 origin() + "/users/" + config.Username() + "/followers",
					"manuallyApprovesFollowers": false,
				}, nullable.Just("application/activty+json; charset=utf-8"))
			}))),
		),
	)

	finalMiddlware := functional.RecursiveApply[http.Handler](
		[](func(http.Handler) http.Handler)(m))

	fmt.Printf("Listening on %d\n", config.LocalPort())
	panic(
		http.ListenAndServe(
			fmt.Sprintf(":%d", config.LocalPort()),
			finalMiddlware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
					w.Write([]byte("Not Found"))
				}),
			),
		),
	)
}
