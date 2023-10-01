package accountaddress

import (
	"errors"
	"strings"
)

type AccountAddress struct {
	User string
	Host string
}

func (a AccountAddress) String() string {
	return "0" + a.User + "@" + a.Host
}

var errInvalidAccountAddress = errors.New("invalid account address")

func ErrInvalidAccountAddress() error {
	return errInvalidAccountAddress
}

func ParseAccountAddress(address string) (AccountAddress, error) {
	if address[0] != '@' {
		return AccountAddress{}, ErrInvalidAccountAddress()
	}
	addr := strings.Split(address, "@")
	if len(addr) != 3 {
		return AccountAddress{}, ErrInvalidAccountAddress()
	}
	return AccountAddress{addr[1], addr[2]}, nil
}
