package posts

import (
	"database/sql"
	"fediverse/application/config"
	"fmt"
	"path"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var lock sync.RWMutex

func CreatePost(body string) error {
	lock.Lock()
	defer lock.Unlock()
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

type Post struct {
	ID          string    `json:"id"`
	WhenCreated time.Time `json:"whenCreated"`
	Body        string    `json:"body"`
}

func GetAllPosts() ([]Post, error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := sql.Open("sqlite3", path.Join(config.OutputDir(), "application.db"))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Query("SELECT id, when_created, body FROM posts")
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for result.Next() {
		var id uint64
		var whenCreated time.Time
		var body string
		if err := result.Scan(&id, &whenCreated, &body); err != nil {
			return nil, err
		}
		posts = append(posts, Post{
			ID:          fmt.Sprintf("%d", id),
			WhenCreated: whenCreated,
			Body:        body,
		})
	}
	return posts, nil
}
