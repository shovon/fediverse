package common

import (
	"fediverse/application/config"
	"fmt"
	"net/url"
)

func Origin() string {
	return config.HttpProtocol() + "://" + config.Hostname()
}

func BaseURL() *url.URL {
	u, err := url.Parse(Origin())
	if err != nil {
		panic(fmt.Errorf("URL of origin %s is not a valid URL. Are we generating base URLs correctly?", Origin()))
	}
	return u
}

func init() {
	BaseURL()
}
