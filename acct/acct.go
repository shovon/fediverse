package acct

import (
	"errors"
	"net/url"
	"strings"
)

// https://datatracker.ietf.org/doc/html/rfc7565

// - Used to identify a user's account at a service provider.

var ErrNotAcct = errors.New("not an acct URI")
var ErrNotValidUserHost = errors.New("not a valid user@host")

type Acct struct {
	User string
	Host string
}

func ParseUserHost(s string) (string, string, error) {
	split := strings.Split(s, "@")
	if len(split) != 2 {
		return "", "", ErrNotValidUserHost
	}
	user, err := url.QueryUnescape(split[0])
	if err != nil {
		return "", "", err
	}
	return user, split[1], nil
}

// ParseAcct parses an acct URI into its components.
func ParseAcct(acct string) (Acct, error) {
	u, err := url.Parse(acct)
	if err != nil {
		return Acct{}, err
	}
	if u.Scheme != "acct" {
		return Acct{}, ErrNotAcct
	}
	if u.Opaque != "" {
		return Acct{}, ErrNotAcct
	}

	user, host, err := ParseUserHost(u.Opaque)
	if err != nil {
		return Acct{}, err
	}

	return Acct{User: user, Host: host}, nil
}
