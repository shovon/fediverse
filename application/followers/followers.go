package followers

import (
	"database/sql"
	"fediverse/application/database"
	"sync"
	"time"
)

var lock sync.RWMutex

// AddFollower adds a following to the database.
//
// Not sure what the implication is for just interpreting the IRI as a string,
// but it will be so much simpler to work with, for now.
func AddFollower(actorIRI string) (int64, error) {
	lock.Lock()
	defer lock.Unlock()
	db, err := database.Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var existingID int64
	err = db.QueryRow("SELECT id FROM followers WHERE actor_iri = ?", actorIRI).Scan(&existingID)

	switch {
	case err == sql.ErrNoRows:
		result, err := db.Exec(
			"INSERT INTO followers (actor_iri) VALUES (?)",
			actorIRI,
		)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	case err != nil:
		return 0, err
	}

	return existingID, nil
}

func RemoveFollower(actorIRI string) error {
	lock.Lock()
	defer lock.Unlock()
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM followers WHERE actor_iri = ?", actorIRI)
	return err

	// TODO: what if the delete silently failed? How are we going to know?
}

type Follower struct {
	ID           string    `json:"id"`
	WhenFollowed time.Time `json:"whenFollowed"`
	ActorIRI     string    `json:"actor_iri"`
}

// GetFollowers gets a following from the database, given an offset.
//
// This function does not do any paging. You have to calculate the offset,
// by multiplying the page number by your desired page size.
func GetFollowers(offset int, limit int) (_ []Follower, err error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Query(
		"SELECT id, when_followed, actor_iri FROM followers ORDER BY when_followed DESC LIMIT ? OFFSET ?",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = result.Close()
	}()
	followings := []Follower{}
	for result.Next() {
		var following Follower
		if err := result.Scan(&following.ID, &following.WhenFollowed, &following.ActorIRI); err != nil {
			return nil, err
		}
		followings = append(followings, following)
	}
	return followings, nil
}
