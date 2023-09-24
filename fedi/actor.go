package fedi

import (
	"fediverse/acct"
	"fediverse/webfinger"
)

type Account struct {
	User string
	Host string
}

func GetActor(account Account) (any, error) {
	// First, we issue a WebFinger request to the specific host, requesting
	// specifically for the application/activity+json content type.

	// But before that, construct a WebFinger URL.

	return webfinger.Lookup(account.Host, acct.Acct(account).String(), []string{"self"})
}
