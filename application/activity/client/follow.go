package client

import "net/http"

// Example follow activity
// {
//   "@context": "https://www.w3.org/ns/activitystreams",
//   "id": "https://mastodon.social/4c303fb1-fca8-4b79-8e8b-a62048e6aa3a",
//   "type": "Follow",
//   "actor": "https://mastodon.social/users/fedipedia",
//   "object": "https://feditest.salrahman.com/activity/users/johndoe"
// }

func Follow(url string) {
	msi := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       "https://mastodon.social/4c303fb1-fca8-4b79-8e8b-a62048e6aa3a",
	}

	req, err := http.NewRequest("POST", url, nil)
}
