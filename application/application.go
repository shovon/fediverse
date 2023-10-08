package application

import (
	"fediverse/acct"
	"fediverse/application/activity/server"
	"fediverse/application/config"
	"fediverse/application/posts"
	"fediverse/functional"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/httplogger"
	"fediverse/jrd"
	"fediverse/json/jsonhttp"
	"fediverse/nodeinfo"
	"fediverse/pathhelpers"
	"fediverse/webfinger"
	"fmt"
	"net/http"
	"net/url"

	"fediverse/application/common"
)

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

func Start() error {
	m := [](func(http.Handler) http.Handler){}

	m = append(m, httplogger.Middleware)

	// TODO: this is really a bad idea. Gotta remove it.
	m = append(m, requestbaseurl.Override(common.Origin()))

	// WebFinger
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

	// NodeInfo
	m = append(m, nodeinfo.CreateNodeInfoMiddleware(common.Origin(), "/nodinfo", func() (nodeinfo.NodeInfoProps, error) {
		count, err := posts.GetPostCount()
		if err != nil {
			return nodeinfo.NodeInfoProps{}, err
		}
		return nodeinfo.NodeInfoProps{
			Software: nodeinfo.SoftwareInfo{
				Name:    "fediverse",
				Version: "0.0.1",
			},
			OpenRegistrations: false,
			Usage: nodeinfo.Usage{
				Users: nodeinfo.UsersStats{
					// TODO: this is where we get a list of all users.
					Total: 1,

					// TODO: actually get all active users in the last 6 months
					ActiveHalfyear: 0,

					// TODO: actually get all activie users in the last 30 days
					ActiveMonth: 0,
				},
				LocalPosts:    count,
				LocalComments: 0, // TODO: actually get all local comments
			},
		}, nil
	}))

	m = append(
		m,
		hh.Processors{
			hh.Method("GET"),
			hh.Route("/useless-api/posts"),
		}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
			return posts.GetAllPosts()
		}))),
	)

	m = append(
		m,
		hh.Processors{
			hh.Method("GET"),
			hh.Route("/useless-api/posts/:id"),
		}.Process(hh.ToMiddleware(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
			return posts.GetPost(hh.GetRouteParam(r, "id"))
		}))),
	)

	m = append(m, server.ActivityPub())

	m = append(m, hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Just an article. Coming soon"))
	})))

	finalMiddlware := functional.RecursiveApply[http.Handler](
		[](func(http.Handler) http.Handler)(m))

	fmt.Printf("Listening on %d\n", config.LocalPort())
	return http.ListenAndServe(
		fmt.Sprintf(":%d", config.LocalPort()),
		finalMiddlware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(404)
				w.Write([]byte("Not Found"))
			}),
		),
	)
}
