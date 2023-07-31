package webfinger

import (
	"encoding/json"
	"fediverse/jrd"
	"fediverse/nullable"
	"net/http"
)

// Notes from rfc7033:
//
// - returns JSON
// - JSON is referred to as a "JRD" (JSON Resource Descriptor)
// - For a person, the type of information that might be discoverable WebFinger
//   includes:
//   - a profile address
//   - identity service
//   - telephone number
//   - preferred avatar
// - For other entities on the Internet, a WebFinger resource might return JRDs
//   containing link relations that enable the client to discover
//   - whether a printer can print colour
//   - the physical location of a server
//   - etc.
// - Information from WebFinger can be for direct human consumption or for
//   systems needing to carry out some function
// - WebFinger is only for static information (e.g. can't be used to get back
//   back dynamic information such as temperature)
// - CORS **MUST** be supported
// - served over a /.well-known/webfinger path
// - must contain a "resource" parameter, which is the URI of the resource
// - may contain one or more "rel" parameters, which is link relation type
// - to request only a subset of the information, the client can specify a
//   "rel" parameter, which whitelists the link relation types that the client
//   is interested in, and the server must filter out

type WebFingerQueryHandler func(string) (jrd.JRD, error)

func CreateHandler(queryHandler WebFingerQueryHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET")
		w.Header().Add("Content-Type", "application/jrd+json")

		// TODO: implement the rel parameter

		j, err := queryHandler(r.URL.Query().Get("resource"))

		rel := r.URL.Query().Get("rel")
		if rel != "" {
			currentLinks, err := j.Links.Value()
			if err != nil {
				links := []jrd.Link{}
				for _, link := range currentLinks {
					if link.Rel == rel {
						links = append(links, link)
					}
				}
				j.Links = nullable.Just(links)
			}
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(j)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(js)
	})
}
