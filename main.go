package main

import (
	"fediverse/acct"
	"fediverse/config"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/nullable"
	"fediverse/webfinger"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func username() string {
	// Just a hard-coded username.
	return "username"
}

func hostname() string {
	return config.Hostname()
}

func httpProtocol() string {
	return config.HttpProtocol()
}

func origin() string {
	return fmt.Sprintf("%s://%s", httpProtocol(), hostname())
}

type UserHost struct {
	Username string
	Host     string
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(webfinger.WebFinger(func(resource string) (jrd.JRD, httperrors.HTTPError) {
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
			user = pathParts[1]
			if len(user) < 2 {
				return jrd.JRD{}, httperrors.BadRequest()
			}
			if len(user) == 2 {
				if user[0] != '@' {
					return jrd.JRD{}, httperrors.BadRequest()
				}
				user = user[1:]
			} else {
				if user[0] != '@' || user[1] != '+' {
					return jrd.JRD{}, httperrors.BadRequest()
				}
				user = user[2:]
			}

		}

		if host != hostname() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		if user != username() {
			return jrd.JRD{}, httperrors.NotFound()
		}

		return jrd.JRD{
			Subject: nullable.Just("acct:" + user + "@" + host),
			Aliases: nullable.Just([]string{
				origin() + "/@" + user,
			}),
			Links: nullable.Just([]jrd.Link{
				{
					Rel:  "http://webfinger.net/rel/profile-page",
					Type: nullable.Just("text/html"),
					Href: origin() + "/@" + user,
				},
				{
					Rel:  "self",
					Type: nullable.Just("application/activity+json"),
					Href: origin() + "/users/" + user,
				},
			}),
		}, nil
	}))
	http.ListenAndServe(":3000", r)
}
