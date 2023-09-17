package database

import (
	"database/sql"
	"fediverse/application/config"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: it's really a bad idea to have configurations in the database package.
func Open() (*sql.DB, error) {
	return sql.Open("sqlite3", path.Join(config.OutputDir(), "application.db"))
}
