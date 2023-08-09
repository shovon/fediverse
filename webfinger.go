package main

import (
	"fediverse/jrd"
	"fediverse/nullable"
	"net/url"
)

func webFingerJRD(userHost UserHost) jrd.JRD {
	user, host := userHost.Username, userHost.Host

	htmlAddress := origin() + "/@" + user
	jsonLDAddress := origin() + "/users/" + host

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
	}
}
