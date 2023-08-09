package main

import (
	"fediverse/acct"
	"fediverse/config"
	"fediverse/functional"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/pathhelpers"
	"fediverse/slices"
	"fediverse/webfinger"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5/middleware"
)

func username() string {
	// Just a hard-coded username.
	return "username"
}

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

func main() {
	m := slices.Pusher[func(http.Handler) http.Handler]{}
	m.Push(middleware.Logger)
	m.Push(webfinger.WebFinger(func(resource string) (jrd.JRD, httperrors.HTTPError) {
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

		if user != username() {
			return jrd.JRD{}, httperrors.NotFound()
		}
		if host != config.Hostname() {
			return jrd.JRD{}, httperrors.NotFound()
		}
		if urlIsValid && urlQuery.Scheme != config.HttpProtocol() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		return webFingerJRD(UserHost{user, host}), nil
	}))
	fmt.Printf("Listening on %d\n", config.LocalPort())
	panic(http.ListenAndServe(fmt.Sprintf(":%d", config.LocalPort()), functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler)(m))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	}))))
}
