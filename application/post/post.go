package post

import (
	"bufio"
	"encoding/json"
	"fediverse/application/config"
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
	f, err := os.OpenFile(path.Join(config.OutputDir()), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	file = f
}

type Post struct {
	WhenCreated time.Time `json:"whenCreated"`
	Body        string    `json:"body"`
}

func CreatePost(body string) error {
	post := Post{
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
