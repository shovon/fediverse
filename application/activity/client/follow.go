package client

import (
	"encoding/json"
	"errors"
)

// Example follow activity
// {
//   "@context": "https://www.w3.org/ns/activitystreams",
//   "id": "https://mastodon.social/4c303fb1-fca8-4b79-8e8b-a62048e6aa3a",
//   "type": "Follow",
//   "actor": "https://mastodon.social/users/fedipedia",
//   "object": "https://feditest.salrahman.com/activity/users/johndoe"
// }

func Follow(objectIRI string, inboxIRI string, followIRI string) error {
	_, err := json.Marshal(map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       "https://mastodon.social/4c303fb1-fca8-4b79-8e8b-a62048e6aa3a",
		"type":     "Follow",
		"actor":    "https://mastodon.social/users/fedipedia",
		"object":   "https://feditest.salrahman.com/activity/users/johndoe",
	})
	if err != nil {
		return err
	}
	return errors.New("not implemented")
	// req, err := http.NewRequest("POST", inboxIRI, strings.NewReader(string(j)))
}
