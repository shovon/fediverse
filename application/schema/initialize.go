package schema

import (
	"database/sql"
	"fediverse/application/config"
	"path"
)

var revisions = []string{
	`
	CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		body TEXT
	);

	CREATE TABLE followers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_followed DATETIME DEFAULT CURRENT_TIMESTAMP,
		follower TEXT
	);

	CREATE TABLE following (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_followed DATETIME DEFAULT CURRENT_TIMESTAMP,
		follower TEXT
	);

	CREATE TABLE inbox (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_received DATETIME DEFAULT CURRENT_TIMESTAMP,

		-- This is a JSON document, more specifically a JSON-LD document.
		body TEXT
	)
	`,
}

func runRevision(db *sql.DB, revision string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := db.Exec(revision); err != nil {
		return err
	}
	if _, err := db.Exec("INSERT INTO schema_revisions (revision) VALUES (?)", revision); err != nil {
		return err
	}
	return tx.Commit()
}

func Initialize() error {
	db, err := sql.Open("sqlite3", path.Join(config.OutputDir(), "application.db"))
	if err != nil {
		return err
	}

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS schema_revisions (revision INTEGER PRIMARY KEY);`,
	); err != nil {
		return err
	}
	result, err := db.Query("SELECT revision FROM schema_revisions ORDER BY revision DESC LIMIT 1;")
	if err != nil {
		return err
	}
	var lastRevision int64
	if err := result.Scan(&lastRevision); err != nil {
		return err
	}
	return nil
}
