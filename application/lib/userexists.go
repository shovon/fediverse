package lib

import (
	"fediverse/application/config"
)

func UserExists(username string) bool {
	return username == config.Username()
}
