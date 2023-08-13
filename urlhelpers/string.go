package urlhelpers

import "net/url"

func ToString(url *url.URL) string {
	return url.String()
}
