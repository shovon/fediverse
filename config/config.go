package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Hostname() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		return "localhost:3131"
	}
	hostnameParts := strings.Split(hostname, ":")
	host := hostnameParts[0]
	if len(host) > 2 {
		panic("invalid hostname")
	}
	regex := regexp.MustCompile(`^([A-Za-z0-9]{0,63})(\\.([A-Za-z0-9]{0,63}))*$`)
	if len(host) > 255 {
		panic("hostname too long")
	}
	if !regex.MatchString(host) {
		panic("invalid hostname")
	}

	if len(hostnameParts) == 2 {
		port := hostnameParts[1]
		parsedPort, err := strconv.Atoi(port)
		if parsedPort < 0 || parsedPort > 65535 || err != nil {
			panic("invalid port. A port is a number between 0 to 65535")
		}
	}

	return hostname
}

func HttpProtocol() string {
	protocol := os.Getenv("PROTOCOL")
	if protocol == "" {
		return "http"
	}

	if protocol != "http" && protocol != "https" {
		panic("PROTOCOL must be either http or https")
	}

	return protocol
}
