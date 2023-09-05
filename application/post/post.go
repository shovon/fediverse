package post

import (
	"database/sql"
	"fediverse/application/config"
	"path"
	"time"
)

type Post struct {
	ID          string    `json:"id"`
	WhenCreated time.Time `json:"whenCreated"`
	Body        string    `json:"body"`
}

func CreatePost(body string) error {
	// TODO: careful with race conditions!
	db, err := sql.Open("sqlite3", path.Join(config.OutputDir(), "application.db"))
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec("INSERT INTO posts (body) VALUES (?)", body); err != nil {
		return err
	}
	return nil
}
