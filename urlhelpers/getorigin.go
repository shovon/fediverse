package urlhelpers

import "net/url"

func GetOrigin(u *url.URL) string {
	newURL := *u
	return newURL.Scheme + "://" + newURL.Host
}
