package schema

import (
	"database/sql"
	"fediverse/application/database"

	_ "github.com/mattn/go-sqlite3"
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
		actor_iri TEXT
	);

	CREATE TABLE following (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_followed DATETIME DEFAULT CURRENT_TIMESTAMP,
		actor_iri TEXT
	);

	CREATE TABLE inbox (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_received DATETIME DEFAULT CURRENT_TIMESTAMP,

		-- This is a JSON document, more specifically a JSON-LD document.
		body TEXT
	)
	`,
	`
	
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
	if _, err := db.Exec("INSERT INTO schema_revisions DEFAULT VALUES", revision); err != nil {
		return err
	}
	return tx.Commit()
}

func runRevisions(db *sql.DB, revisions []string) error {
	for _, revision := range revisions {
		if err := runRevision(db, revision); err != nil {
			return err
		}
	}
	return nil
}

func Initialize() (err error) {
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS schema_revisions (revision INTEGER PRIMARY KEY AUTOINCREMENT);`,
	); err != nil {
		return err
	}
	result, err := db.Query("SELECT revision FROM schema_revisions ORDER BY revision DESC LIMIT 1;")
	if err != nil {
		return err
	}
	defer func() {
		err = result.Close()
	}()
	hasNext := result.Next()
	if !hasNext {
		if err := runRevisions(db, revisions); err != nil {
			return err
		}
	} else {
		var lastRevision int64
		if err := result.Scan(&lastRevision); err != nil {
			return err
		}
		if lastRevision >= int64(len(revisions)) {
			return nil
		}
		if err := runRevisions(db, revisions[lastRevision:]); err != nil {
			return err
		}
	}
	return nil
}
