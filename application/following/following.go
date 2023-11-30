package following

import (
	"database/sql"
	"fediverse/accountaddress"
	"fediverse/application/database"
	"sync"
	"time"
)

var lock sync.RWMutex

type Following struct {
	ID           string    `json:"id"`
	WhenFollowed time.Time `json:"whenFollowed"`
	ActorIRI     string    `json:"actorIri"`
}

// TODO: handle a way to update the current user IRIs

// TODO: all of these functions should also factor in the follower/followee.
//
// Right now, all of these functions only really care about the followee.

// GetFollowing gets a following from the database, given an offset.
//
// This function does not do any paging. You have to calculate the offset,
// by multiplying the page number by your desired page size.
func GetFollowing(offset int, limit int) (_ []Following, err error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Query(
		"SELECT id, when_followed, actor_iri FROM following ORDER BY when_followed DESC LIMIT ? OFFSET ?",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeError := result.Close()
		if closeError != nil {
			err = closeError
		}
	}()
	followings := []Following{}
	for result.Next() {
		var following Following
		if err := result.Scan(
			&following.ID,
			&following.WhenFollowed,
			&following.ActorIRI,
		); err != nil {
			return nil, err
		}
		followings = append(followings, following)
	}
	return followings, nil
}

// TODO: this should also really factor in a follower/followee relationship
//   somehow

func GetSingleFollowingID(actorID string) (Following, error) {
	lock.RLock()
	defer lock.RUnlock()
	db, err := database.Open()
	if err != nil {
		return Following{}, err
	}
	defer db.Close()
	var following Following
	err = db.QueryRow(
		"SELECT id, when_followed, actor_iri FROM following WHERE id = ?",
		actorID,
	).Scan(
		&following.ID,
		&following.WhenFollowed,
		&following.ActorIRI,
	)
	return following, err
}

// AddFollowing adds a following to the database.
//
// Not sure what the implication is for just interpreting the IRI as a string,
// but it will be so much simpler to work with, for now.
func AddFollowing(actorIRI string) (int64, error) {
	lock.Lock()
	defer lock.Unlock()
	db, err := database.Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var existingID int64
	err = db.QueryRow("SELECT id FROM following WHERE actor_iri = ?", actorIRI).Scan(&existingID)

	switch {
	case err == sql.ErrNoRows:
		result, err := db.Exec(
			"INSERT INTO following (actor_iri) VALUES (?)",
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

func RemoveFollowing(actorIRI string) error {
	lock.Lock()
	defer lock.Unlock()
	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM following WHERE actor_iri = ?", actorIRI)
	return err
}

func AcknowledgeFollowing(id int) {
	lock.Lock()
	defer lock.Unlock()

	db, err := database.Open()
	if err != nil {
		return
	}
	defer db.Close()

	if _, err := db.Exec(
		"UPDATE following SET has_accepted_follow_request = 1 WHERE id = ?", id,
	); err != nil {
		return
	}
}

func FollowRequestAccepted(address accountaddress.AccountAddress) error {
	lock.Lock()
	defer lock.Unlock()

	db, err := database.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(
		"UPDATE following SET accepted = 1 WHERE account_address_user = ? AND account_address_host = ?", address.User, address.Host,
	); err != nil {
		return err
	}

	return nil
}
