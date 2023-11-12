package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fediverse/application/apphttp"
	"fediverse/httphelpers/httpsig/cavage"
	"fediverse/httphelpers/rfc3230"
	"fediverse/nullable"
	"fediverse/security/rsassapkcsv115sha256"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Example follow activity
// {
//   "@context": "https://www.w3.org/ns/activitystreams",
//   "id": "https://mastodon.social/4c303fb1-fca8-4b79-8e8b-a62048e6aa3a",
//   "type": "Follow",
//   "actor": "https://mastodon.social/users/fedipedia",
//   "object": "https://feditest.salrahman.com/activity/users/johndoe"
// }

type FollowActivityIRI string
type AcceptActivityIRI string
type SigningKeyIRI string
type SenderIRI string
type ObjectIRI string
type InboxURL string

func Follow(
	signingKey *rsa.PrivateKey,
	signingKeyIRI SigningKeyIRI,
	followActivityIRI FollowActivityIRI,
	senderIRI SenderIRI,
	recipientID ObjectIRI,
	inboxURL InboxURL,
) error {
	b, err := json.Marshal(map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       followActivityIRI,
		"type":     "Follow",
		"actor":    senderIRI,
		"object":   recipientID,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", string(inboxURL), bytes.NewReader(b))
	if err != nil {
		return err
	}

	digesters := []rfc3230.Digester{apphttp.SHA256Digest{}}
	if err := rfc3230.AddDigestToRequest(req, digesters); err != nil {
		return err
	}

	fmt.Println(req.Header.Get("Digest"))

	signer := rsassapkcsv115sha256.Base64Signer(signingKey)

	if err := cavage.AddSignatureToRequest(req, cavage.Params{
		KeyID:     nullable.Just(string(signingKeyIRI)),
		Algorithm: nullable.Just("hs2019"),
		Headers:   nullable.Just([]string{"digest", "(created)", "(request-target)"}),
		Created:   time.Now(),
		// TODO: consider adding an expries field
	}, signer); err != nil {
		return err
	}

	fmt.Println(req.Header.Get("Signature"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// We want to send an activity of the following form:
//
// {
//   "@context": "https://www.w3.org/ns/activitystreams",
//   "id": "https://techhub.social/users/manlycoffee#accepts/follows/1129830",
//   "type": "Accept",
//   "actor": "https://techhub.social/users/manlycoffee",
//   "object": {
//     "id": "https://feditest.salrahman.com/activity/actors/johndoe/following/1",
//     "type": "Follow",
//     "actor": "https://feditest.salrahman.com/activity/actors/johndoe",
//     "object": "https://techhub.social/users/manlycoffee"
//   }
// }

func AcceptFollow(
	signingKey *rsa.PrivateKey,
	signingKeyIRI SigningKeyIRI,
	acceptActivityIRI AcceptActivityIRI,
	senderIRI SenderIRI,
	object any,
	inboxURL InboxURL,
) error {
	b, err := json.Marshal(map[string]any{
		"@context": "https://www.w3.org/ns/activitystreams",
		"id":       acceptActivityIRI,
		"type":     "Accept",
		"actor":    senderIRI,
		"object":   object,
	})
	if err != nil {
		return err
	}

	fmt.Println("Follow activity payload:", string(b))

	req, err := http.NewRequest("POST", string(inboxURL), bytes.NewReader(b))
	if err != nil {
		return err
	}

	digesters := []rfc3230.Digester{apphttp.SHA256Digest{}}
	if err := rfc3230.AddDigestToRequest(req, digesters); err != nil {
		return err
	}

	fmt.Println(req.Header.Get("Digest"))

	signer := rsassapkcsv115sha256.Base64Signer(signingKey)

	if err := cavage.AddSignatureToRequest(req, cavage.Params{
		KeyID:     nullable.Just(string(signingKeyIRI)),
		Algorithm: nullable.Just("hs2019"),
		Headers:   nullable.Just([]string{"digest", "(created)", "(request-target)"}),
		Created:   time.Now(),
		// TODO: consider adding an expries field
	}, signer); err != nil {
		return err
	}

	fmt.Println(req.Header.Get("Signature"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
