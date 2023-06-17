package acct

import (
	"errors"
	"strings"
)

// - Used to identify a user's account at a service provider.

var ErrNotAcct = errors.New("not an acct URI")

type Acct struct {
	User string
	Host string
}

// ParseAcct parses an acct URI into its components.
func ParseAcct(acct string) (Acct, error) {
	if !strings.HasPrefix(acct, "acct:") {
		return Acct{}, ErrNotAcct
	}

	split := strings.Split(acct, ":")
	if len(split) != 2 {
		return Acct{}, ErrNotAcct
	}

	split = strings.Split(split[1], "@")
	if len(split) != 2 {
		return Acct{}, ErrNotAcct
	}

	return Acct{
		User: split[0],
		Host: split[1],
	}, nil
}
