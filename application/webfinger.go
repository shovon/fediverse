package application

import (
	"fediverse/application/activity/routes"
	"fediverse/application/common"
	"fediverse/jrd"
	"fediverse/nullable"
	"net/url"
)

func webFingerJRD(userHost UserHost) jrd.JRD {
	user, host := userHost.Username, userHost.Host

	htmlAddress := common.Origin() + "/@" + user

	// TODO: ideally we should be soft-coding the "/activity/users" part
	jsonLDAddress := common.Origin() + routes.Activity{}.Actors().Actor().Route().FullRoute(user)

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
