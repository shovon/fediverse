package main

import (
	"fediverse/acct"
	"fediverse/config"
	"fediverse/functional"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/nullable"
	"fediverse/slices"
	"fediverse/webfinger"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

func username() string {
	// Just a hard-coded username.
	return "username"
}

func origin() string {
	return fmt.Sprintf("%s://%s", config.HttpProtocol(), config.Hostname())
}

type UserHost struct {
	Username string
	Host     string
}

func main() {
	m := slices.Pusher[func(http.Handler) http.Handler]{}
	m.Push(middleware.Logger)
	m.Push(webfinger.WebFinger(func(resource string) (jrd.JRD, httperrors.HTTPError) {
		acctQuery, errAcct := acct.ParseAcct(resource)
		urlQuery, errURL := url.Parse(resource)

		if errAcct != nil && errURL != nil {
			return jrd.JRD{}, httperrors.BadRequest()
		}

		var user, host string

		if errAcct == nil {
			user = acctQuery.User
			host = acctQuery.Host
		} else {
			pathParts := strings.Split(urlQuery.Path, "/")
			if len(pathParts) != 2 && len(pathParts) != 3 {
				return jrd.JRD{}, httperrors.BadRequest()
			}

			host = urlQuery.Host

			if len(pathParts) == 2 {
				user = pathParts[1]
				user = user[1:]
			} else {
				if pathParts[1] != "users" {
					return jrd.JRD{}, httperrors.BadRequest()
				}
				user = pathParts[2]
			}
		}

		if host != config.Hostname() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		if user != username() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		htmlAddress := origin() + "/@" + user
		jsonLDAddress := origin() + "/users/" + user

		return jrd.JRD{
			Subject: nullable.Just("acct:" + user + "@" + url.QueryEscape(host)),
			Aliases: nullable.Just([]string{
				htmlAddress,
				jsonLDAddress,
			}),
			Links: nullable.Just([]jrd.Link{
				{
					Rel:  "http://webfinger.net/rel/profile-page",
					Type: nullable.Just("text/html"),
					Href: htmlAddress,
				},
				{
					Rel:  "self",
					Type: nullable.Just("application/activity+json"),
					Href: jsonLDAddress,
				},
			}),
		}, nil
	}))
	fmt.Printf("Listening on %d\n", config.LocalPort())
	panic(http.ListenAndServe(fmt.Sprintf(":%d", config.LocalPort()), functional.RecursiveApply[http.Handler]([](func(http.Handler) http.Handler)(m))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	}))))
}
