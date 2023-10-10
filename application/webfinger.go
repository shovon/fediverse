package application

import (
	"fediverse/application/activity/server"
	"fediverse/application/common"
	"fediverse/jrd"
	"fediverse/nullable"
	"fediverse/pathhelpers"
	"net/url"
)

func webFingerJRD(userHost UserHost) jrd.JRD {
	user, host := userHost.Username, userHost.Host
	htmlAddress := common.Origin() + "/@" + user
	jsonLDAddress := common.Origin() + pathhelpers.FillFields(server.UserRoute, map[string]string{"username": user})

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
