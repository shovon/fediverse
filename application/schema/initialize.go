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
		
		actor_iri TEXT,

		UNIQUE(actor_iri)
	);

	CREATE INDEX followers_actor_iri_index ON followers(actor_iri);

	CREATE TABLE following (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_followed DATETIME DEFAULT CURRENT_TIMESTAMP,

		actor_iri TEXT,

		has_accepted_follow_request INT DEFAULT 0,

		UNIQUE(actor_iri)
	);

	CREATE INDEX following_actor_iri_index ON following(actor_iri);

	CREATE TABLE inbox (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		when_received DATETIME DEFAULT CURRENT_TIMESTAMP,

		-- This is a JSON document, more specifically a JSON-LD document.
		body TEXT
	);
	`,

	// In the future, we may want to add a cache to hold the actor's IRI
}

func runRevision(db *sql.DB, revision string) (err error) {
	err = nil
	tx, err := db.Begin()
	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = e
			}
		} else {
			if e := tx.Commit(); e != nil {
				err = e
			}
		}
	}()
	if err != nil {
		return
	}
	if _, err = db.Exec(revision); err != nil {
		return
	}
	if _, err = db.Exec("INSERT INTO schema_revisions DEFAULT VALUES", revision); err != nil {
		return
	}
	return
}

func runRevisions(db *sql.DB, revisions []string) error {
	for _, revision := range revisions {
		if err := runRevision(db, revision); err != nil {
			return err
		}
	}
	return nil
}

func Initialize() error {
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
	result, err = db.Query("SELECT revision FROM schema_revisions")
	if err != nil {
		return err
	}
	return err
}
