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

func (a Acct) String() string {
	return "acct:" + url.QueryEscape(a.User) + "@" + url.QueryEscape(a.Host)
}

// TODO: unit test this

func ParseUserHost(s string) (Acct, error) {
	split := strings.Split(s, "@")
	if len(split) != 2 {
		return Acct{}, ErrNotValidUserHost
	}
	user, err := url.QueryUnescape(split[0])
	if err != nil {
		return Acct{}, err
	}
	host, err := url.QueryUnescape(split[1])
	if err != nil {
		return Acct{}, err
	}
	return Acct{User: user, Host: host}, nil
}

// TODO: unit test this

// ParseAcct parses an acct URI into its components.
func ParseAcct(acct string) (Acct, error) {
	u, err := url.Parse(acct)
	if err != nil {
		return Acct{}, err
	}
	if u.Scheme != "acct" {
		return Acct{}, ErrNotAcct
	}
	if u.Opaque == "" {
		return Acct{}, ErrNotAcct
	}

	a, err := ParseUserHost(u.Opaque)
	if err != nil {
		return Acct{}, err
	}

	return Acct{User: a.User, Host: a.Host}, nil
}
