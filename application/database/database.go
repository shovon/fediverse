package database

import (
	"database/sql"
	"fediverse/application/config"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func Open() (*sql.DB, error) {
	return sql.Open("sqlite3", path.Join(config.OutputDir(), "application.db"))
}
