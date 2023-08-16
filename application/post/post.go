package post

import (
	"bufio"
	"encoding/json"
	"fediverse/application/config"
	"fediverse/application/id"
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

var writeMutext sync.Mutex

var file *os.File

func init() {
	os.MkdirAll(config.OutputDir(), os.ModePerm)
	f, err := os.OpenFile(path.Join(config.OutputDir(), "posts.jsonl"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	file = f
}

type Post struct {
	ID          string    `json:"id"`
	WhenCreated time.Time `json:"whenCreated"`
	Body        string    `json:"body"`
}

func CreatePost(body string) error {
	i, err := id.Generate()
	if err != nil {
		return err
	}

	post := Post{
		ID:          i,
		WhenCreated: time.Now(),
		Body:        body,
	}
	b, err := json.Marshal(post)
	if err != nil {
		return err
	}
	writeMutext.Lock()
	defer writeMutext.Unlock()

	writer := bufio.NewWriter(file)

	_, err = fmt.Fprintln(writer, string(b))
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
