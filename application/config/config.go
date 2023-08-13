package config

import (
	"fediverse/pathhelpers"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var username string
var displayName string
var hostname string
var httpProtocol string
var localPort uint16
var outputDir string

func getUsername() {
	username = os.Getenv("USERNAME")
	if username == "" {
		panic("USERNAME environment variable is not set")
	}
}

func getDisplayName() {
	displayName = os.Getenv("DISPLAY_NAME")
	if displayName == "" {
		displayName = username
	}
}

func getLocalPort() {
	port := os.Getenv("LOCAL_PORT")
	if port == "" {
		localPort = 3131
		return
	}

	parsedPort, err := strconv.Atoi(port)
	if parsedPort < 0 || parsedPort > 65535 || err != nil {
		panic("invalid port. A port is a number between 0 to 65535")
	}
	localPort = uint16(parsedPort)
}

func getHostname() {
	// TODO: the hostname is a lot more than just a FQDN and a port number. This
	// thing should handle more use cases.

	hostname = os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = "localhost:3131"
		return
	}
	hostnameParts := strings.Split(hostname, ":")
	host := hostnameParts[0]
	if len(hostnameParts) > 2 {
		panic(fmt.Sprintf("invalid hostname %s. Way too many colons. Total number of colons %d", hostname, len(hostnameParts)))
	}
	regex := regexp.MustCompile(`^([A-Za-z0-9]{0,63})(\\.([A-Za-z0-9]{0,63}))*$`)
	if len(host) > 255 {
		panic("hostname too long")
	}
	if !regex.MatchString(host) {
		panic("invalid hostname " + hostname)
	}

	if len(hostnameParts) == 2 {
		port := hostnameParts[1]
		parsedPort, err := strconv.Atoi(port)
		if parsedPort < 0 || parsedPort > 65535 || err != nil {
			panic("invalid port. A port is a number between 0 to 65535")
		}
	}
}

func getHTTPProtocol() {
	httpProtocol = os.Getenv("HTTP_PROTOCOL")
	if httpProtocol == "" {
		httpProtocol = "http"
		return
	}

	if httpProtocol != "http" && httpProtocol != "https" {
		panic("PROTOCOL must be either http or https")
	}
}

func getOutputDir() {
	odir := os.Getenv("OUTPUT_DIR")
	if odir == "" {
		panic("an output directory must be specified")
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	outputDir = pathhelpers.Resolve(wd, odir)
}

func init() {
	getUsername()
	getDisplayName()
	getLocalPort()
	getHostname()
	getHTTPProtocol()
	getOutputDir()
}

func Username() string {
	return username
}

func DisplayName() string {
	return displayName
}

func Hostname() string {
	return hostname
}

func HttpProtocol() string {
	return httpProtocol
}

func LocalPort() uint16 {
	return localPort
}

func OutputDir() string {
	return outputDir
}
