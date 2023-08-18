package requestbaseurl

import (
	"net/http"
	"net/url"
)

func GetRequestUrl(r *http.Request) (*url.URL, error) {
	if u, ok := r.Context().Value(overriden{}).(*url.URL); ok {
		return u, nil
	}

	u, err := r.URL.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}
	u.Host = r.Host
	if r.TLS != nil {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	return u, err
}
