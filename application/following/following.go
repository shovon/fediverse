package following

import (
	"fediverse/application/database"
	"sync"
	"time"
)

var lock sync.RWMutex

// AddFollwing adds a following to the database.
//
// Not sure what the implication is for just interpreting the IRI as a string,
// but it will be so much simpler to work with, for now.
func AddFollwing(i string) error {
	lock.Lock()
	defer lock.Unlock()
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec("INSERT INTO following (following) VALUES (?)", i); err != nil {
		return err
	}
	return nil
}

type Following struct {
	ID           string    `json:"id"`
	WhenFollowed time.Time `json:"whenFollowed"`
	ActorIRI     string    `json:"actor_iri"`
}

// GetFollowing gets a following from the database, given an offset.
//
// This function does not do any paging. You have to calculate the offset,
// by multiplying the page number by your desired page size.
func GetFollowing(offset int, limit int) ([]Following, error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Query("SELECT id, when_followed, actor_iri FROM following ORDER BY when_followed DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	followings := []Following{}
	for result.Next() {
		var following Following
		if err := result.Scan(&following.ID, &following.WhenFollowed, &following.ActorIRI); err != nil {
			return nil, err
		}
		followings = append(followings, following)
	}
	return followings, nil
}
