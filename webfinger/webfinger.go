package webfinger

import (
	"encoding/json"
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/jrd/jrdhttp"
	"fediverse/wellknown"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type WebFingerQueryHandler func(string) (jrd.JRD, httperrors.HTTPError)

// CreateHandler creates a handler for the WebFinger endpoint.
//
// As far as the API exposed by the `CreateHandler` endpoint is concerned, it is
// a somewhat opininated API, as well as the implementation itself is not fully
// compliant. One example is that the handler does not filter out requests that
// aren't done via HTTPS.
//
// One way that this implemenation is opinionated is that if the entire
// request/response cycle yields an error, what to respond with in the body is
// not defined by the WebFinger specification, but this implementation has opted
// instead to respond with a status code and an empty body.
func WebFinger(queryHandler WebFingerQueryHandler) func(http.Handler) http.Handler {
	return httphelpers.Method("GET").Process(wellknown.WellKnown(
		"webfinger",
		CORS(jrdhttp.CreateJRDHandler(func(r *http.Request) (jrd.JRD, httperrors.HTTPError) {
			j, err := queryHandler(r.URL.Query().Get("resource"))

			if err != nil {
				return j, err
			}

			if subject, ok := j.Subject.Value(); !ok || subject == "" {
				return j, httperrors.InternalServerError()
			}

			j = HandleRel(j, r)
			return j, nil
		})),
	))
}

// Lookup issues an HTTP request to grab the JRD for the given resource.
func Lookup(host string, resource string, rel []string) (jrd.JRD, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s", host))
	if err != nil {
		return jrd.JRD{}, err
	}
	u.Path = "/.well-known/webfinger"
	q := u.Query()
	q.Set("resource", resource)
	for _, r := range rel {
		q.Add("rel", r)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return jrd.JRD{}, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return jrd.JRD{}, err
	}

	// TODO: look at response code.

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return jrd.JRD{}, err
	}

	var j jrd.JRD
	err = json.Unmarshal(b, &j)
	if err != nil {
		return jrd.JRD{}, err
	}
	return j, nil
}
