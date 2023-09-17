package posts

import (
	"fediverse/application/database"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var lock sync.RWMutex

func CreatePost(body string) (err error) {
	lock.Lock()
	defer lock.Unlock()
	db, err := database.Open()
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

func GetPost(index string) (Post, error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return Post{}, err
	}
	defer db.Close()
	result, err := db.Query("SELECT id, when_created, body FROM posts WHERE id = ?", index)
	if err != nil {
		return Post{}, err
	}
	defer func() {
		err = result.Close()
	}()
	if !result.Next() {
		return Post{}, fmt.Errorf("no such post")
	}
	var id string
	var whenCreated time.Time
	var body string
	if err := result.Scan(&id, &whenCreated, &body); err != nil {
		return Post{}, err
	}
	return Post{
		ID:          id,
		WhenCreated: whenCreated,
		Body:        body,
	}, nil
}

func GetAllPosts() (_ []Post, err error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Query("SELECT id, when_created, body FROM posts")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = result.Close()
	}()
	posts := []Post{}
	for result.Next() {
		var id string
		var whenCreated time.Time
		var body string
		if err := result.Scan(&id, &whenCreated, &body); err != nil {
			return nil, err
		}
		posts = append(posts, Post{
			ID:          id,
			WhenCreated: whenCreated,
			Body:        body,
		})
	}
	return posts, nil
}

func GetPostCount() (_ uint64, err error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	result, err := db.Query("SELECT COUNT(*) FROM posts")
	if err != nil {
		return 0, err
	}
	defer func() {
		err = result.Close()
	}()
	if !result.Next() {
		return 0, fmt.Errorf("fatal error")
	}
	var count uint64
	if err := result.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
