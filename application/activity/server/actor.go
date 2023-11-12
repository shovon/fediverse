package server

import (
	"encoding/json"
	activityclient "fediverse/application/activity/client"
	"fediverse/application/activity/server/orderedcollection"
	"fediverse/application/config"
	"fediverse/application/followers"
	"fediverse/application/following"
	"fediverse/application/keymanager"
	"fediverse/application/lib"
	hh "fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/httphelpers/requestbaseurl"
	"fediverse/json/jsonhttp"
	"fediverse/jsonldhelpers"
	"fediverse/pathhelpers"
	"fediverse/security/rsahelpers"
	"fediverse/slices"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/piprate/json-gold/ld"
)

type Following string
type Follower string

func searchUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lib.UserExists(hh.GetRouteParam(r, "username")) {
			httperrors.NotFound().ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func actor() func(http.Handler) http.Handler {
	return hh.ApplyMiddlewares(hh.MiddlewaresList{
		// The main user route
		hh.Processors{
			hh.Method("GET"),
			hh.Route(UserRoute),
		}.Process(hh.ToMiddleware(searchUser(jsonhttp.JSONResponder(func(r *http.Request) (any, error) {
			key := keymanager.GetPrivateKey()

			// TODO: this should ideally be cached.
			pubKeyString, err := rsahelpers.PublicKeyToPKIXString(&key.PublicKey)
			if err != nil {
				return nil, httperrors.InternalServerError()
			}

			origin := requestbaseurl.GetRequestOrigin(r)

			params := map[string]string{
				// TODO: this should be soft-coded.
				"username": hh.GetRouteParam(r, "username"),
			}

			actorRoot := origin + pathhelpers.FillFields(UserRoute, params)

			return map[string]any{
				"@context": []interface{}{
					"https://www.w3.org/ns/activitystreams",
					"https://w3id.org/security/v1",
				},

				"id": actorRoot,

				"type": "Person",

				// TODO: this should be soft-coded. That is, retrieve the username,
				//   given some lookup invocation.
				"preferredUsername": config.Username(),

				// TODO: this should be soft-coded. That is, retrieve the display name,
				//   given some lookup invocation.
				"name": config.DisplayName(),

				// TODO: also find a way to soft code this.
				"summary": "<p>This person doesn't have a bio yet.</p>",

				"following": origin + pathhelpers.FillFields(FollowingRoute, params),
				"followers": origin + pathhelpers.FillFields(FollowersRoute, params),
				"inbox":     origin + pathhelpers.FillFields(InboxRoute, params),
				"outbox":    origin + pathhelpers.FillFields(OutboxRoute, params),
				"liked":     origin + pathhelpers.FillFields(LikedRoute, params),

				// TODO: manually approving followers is definitely an important
				//   feature.
				"manuallyApprovesFollowers": false,
				"publicKey": map[string]any{
					"id":           actorRoot + "#main-key",
					"owner":        actorRoot,
					"publicKeyPem": pubKeyString,
				},

				// TODO:
				"endpoints": map[string]any{
					"sharedInbox": origin + SharedInbox,
				},
			}, nil
		})))),

		// The followers collection.
		hh.Processors{
			hh.Route(FollowingRoute),
		}.Process(hh.ApplyMiddlewares(hh.MiddlewaresList{
			searchUser,
			orderedcollection.Middleware(
				orderedcollection.NewOrderedCollection[Following](
					func(hh.ReadOnlyRequest) uint64 {
						return 0
					},
					func(hh.ReadOnlyRequest, orderedcollection.ItemsFunctionParams) []Following {
						return []Following{}
					},
				),
			),
		})),

		// The following collection
		hh.Processors{
			hh.Route(FollowersRoute),
		}.Process(hh.ApplyMiddlewares(hh.MiddlewaresList{
			searchUser,
			orderedcollection.Middleware(
				orderedcollection.NewOrderedCollection[Following](
					func(hh.ReadOnlyRequest) uint64 {
						return 0
					},
					func(hh.ReadOnlyRequest, orderedcollection.ItemsFunctionParams) []Following {
						return []Following{}
					},
				),
			),
		})),

		// The inbox route.
		hh.Processors{
			hh.Route(InboxRoute),
		}.Process(hh.ToMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Inbox")

			d, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Failed to read body")
				w.WriteHeader(500)
				return
			}
			var parsedActivity map[string]any
			if err := json.Unmarshal(d, &parsedActivity); err != nil {
				fmt.Println("Failed to unmarshal JSON")
				w.WriteHeader(400)
				return
			}

			proc := ld.NewJsonLdProcessor()
			options := ld.NewJsonLdOptions("")

			// TODO: add a fallback for the event that the context is not provided.
			//   or is invalid.
			expanded, err := proc.Expand(parsedActivity, options)
			if err != nil {
				fmt.Println("Failed to expand JSON-LD")
				w.WriteHeader(400)
				return
			}

			if len(expanded) != 1 {
				fmt.Println("Expected exactly one JSON-LD document")
				w.WriteHeader(400)
				return
			}

			// Parse tthe activity.
			activity, ok := slices.First(expanded)
			if !ok {
				fmt.Println("Failed to determine activity")
				w.WriteHeader(400)
				return
			}

			switch {
			// Accept activity
			case jsonldhelpers.IsType(activity, "https://www.w3.org/ns/activitystreams#Accept"):
				doc, ok := activity.(map[string]any)
				if !ok {
					fmt.Println("Failed to cast activity to map")
					w.WriteHeader(500)
					return
				}
				obj, ok := doc["https://www.w3.org/ns/activitystreams#object"]
				if !ok {
					fmt.Println("Unable to determine object of Accept activity")
					w.WriteHeader(400)
					return
				}
				if !jsonldhelpers.IsType(obj, "https://www.w3.org/ns/activitystreams#Follow") {
					fmt.Println("Unknown activity to 'accept'")
					w.WriteHeader(400)
					return
				}

				id, ok := jsonldhelpers.GetNodeID(obj)
				if !ok {
					fmt.Println("Unable to determine ID of Follow activity")
					w.WriteHeader(400)
					return
				}

				components := strings.Split(id, "/")
				if len(components) == 0 {
					fmt.Println("Invalid ID string supplied")
					w.WriteHeader(400)
					return
				}

				followID := components[len(components)-1]
				i, err := strconv.Atoi(followID)
				if err != nil {
					fmt.Println("Unable to determine following ID")
					w.WriteHeader(500)
					return
				}
				following.AcknowledgeFollowing(i)

			// Follow activity.
			case jsonldhelpers.IsType(activity, "https://www.w3.org/ns/activitystreams#Follow"):
				// {
				//   "@context":"https://www.w3.org/ns/activitystreams",
				//   "id":"https://techhub.social/a1456ff0-ca04-4c0c-83b8-38df5c693f85",
				//   "type":"Follow",
				//   "actor":"https://techhub.social/users/manlycoffee",
				//   "object":"https://feditest.salrahman.com/activity/actors/john10"
				// }

				doc, ok := activity.(map[string]any)

				// Step 1: grab the object of the body.
				if !ok {
					fmt.Fprint(os.Stderr, "Failed to cast activity to map")
					w.WriteHeader(400)
					return
				}

				actor := jsonldhelpers.GetObjects(doc, "https://www.w3.org/ns/activitystreams#actor")
				if actor == nil {
					fmt.Fprintln(os.Stderr, "There does not appear to be an actor associated with the user")
					w.WriteHeader(400)
					return
				}

				if len(actor) != 1 {
					fmt.Fprintf(os.Stderr, "Expected only a single actor but got %d actors\n", len(actor))
					w.WriteHeader(400)
					return
				}

				firstActor, ok := slices.First(actor)
				if !ok {
					fmt.Fprintln(os.Stderr, "Unable to determine actor")
					w.WriteHeader(400)
					return
				}

				actorIRI, ok := jsonldhelpers.GetNodeID(firstActor)
				if !ok {
					fmt.Fprintln(os.Stderr, "Unable to determine actor IRI")
					w.WriteHeader(400)
					return
				}

				objectObject, ok := slices.First(jsonldhelpers.GetObjects(doc, "https://www.w3.org/ns/activitystreams#object"))
				if !ok {
					fmt.Fprintln(os.Stderr, "Unable to determine object")
					w.WriteHeader(400)
					return
				}

				objectIRI, ok := jsonldhelpers.GetNodeID(objectObject)
				if !ok {
					fmt.Fprintln(os.Stderr, "Unable to determine object IRI")
					w.WriteHeader(400)
					return
				}

				origin := requestbaseurl.GetRequestOrigin(r)

				params := map[string]string{
					// TODO: this should be soft-coded.
					"username": hh.GetRouteParam(r, "username"),
				}

				fmt.Println("Origin", origin)
				actorRoot := origin + pathhelpers.FillFields(UserRoute, params)

				if objectIRI != actorRoot {
					fmt.Fprintln(os.Stderr, "The object IRI must be the same as the actor IRI")
					w.WriteHeader(400)
					return
				}

				// TODO: this should really be placed at the bottom.
				followerID, err := followers.AddFollower(actorIRI)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to add follower %e", err)
					w.WriteHeader(500)
					return
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

				// acceptDoc := map[string]any{
				// 	"@context": "https://www.w3.org/ns/activitystreams",
				// 	"id":       actorRoot + "#accepts/follows/" + strconv.FormatInt(followerID, 10),
				// 	"type":     "Accept",
				// 	"actor":    actorRoot,
				// 	"object":   doc,
				// }

				req, err := http.NewRequest("GET", actorIRI, nil)
				req.Header.Add("Accept", "application/activity+json")

				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to create request to %s: %e", actorIRI, err)
					w.WriteHeader(500)
					return
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to send request to %s: %e", actorIRI, err)
					w.WriteHeader(500)
					return
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to read response body: %s\n", err.Error())
					w.WriteHeader(400)
					return
				}

				var parsed any
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to unmarshal JSON: %s\n", err.Error())
					w.WriteHeader(400)
					return
				}

				objects, err := proc.Expand(parsed, options)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to expand JSON-LD: %s\n", err.Error())
					w.WriteHeader(400)
					return
				}

				if len(objects) != 1 {
					fmt.Fprintf(os.Stderr, "Expected exactly one JSON-LD document but got %d\n", len(objects))
					w.WriteHeader(400)
					return
				}

				object, ok := slices.First(objects)
				if !ok {
					fmt.Fprintln(os.Stderr, "Unable to determine object")
					w.WriteHeader(400)
					return
				}

				inbox, ok := jsonldhelpers.GetIDFromPredicate(object, "http://www.w3.org/ns/ldp#inbox")
				if !ok {
					fmt.Fprintln(os.Stderr, "Unable to determine inbox")
					w.WriteHeader(400)
					return
				}

				privateKey := keymanager.GetPrivateKey()
				signingKeyIRI := actorRoot + "#main-key"
				acceptActivityIRI := actorRoot + "#accepts/follows/" + strconv.FormatInt(followerID, 10)
				senderIRI := actorRoot
				obj := parsedActivity
				inboxURL := inbox

				fmt.Println(acceptActivityIRI)

				err = activityclient.AcceptFollow(
					privateKey,
					activityclient.SigningKeyIRI(signingKeyIRI),
					activityclient.AcceptActivityIRI(acceptActivityIRI),
					activityclient.SenderIRI(senderIRI),
					obj,
					activityclient.InboxURL(inboxURL),
				)

				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to send Accept activity: %s\n", err.Error())
					w.WriteHeader(500)
					return
				}

			default:
				fmt.Println("Unknown activity type")
				fmt.Println(string(d))
			}
		}))),
	})
}
