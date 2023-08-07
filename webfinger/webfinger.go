package webfinger

import (
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/jrd/jrdhttp"
	"fediverse/wellknown"
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
	m := httphelpers.Middlewares{}
	m.Use(CORS)
	m.Use(httphelpers.ToMiddleware(jrdhttp.CreateJRDHandler(func(r *http.Request) (jrd.JRD, httperrors.HTTPError) {
		j, err := queryHandler(r.URL.Query().Get("resource"))

		if err != nil {
			return j, err
		}

		{
			subject, err := j.Subject.Value()
			if err != nil || subject == "" {
				return j, httperrors.InternalServerError()
			}
		}

		j = HandleRel(j, r)
		return j, nil
	})))
	return wellknown.WellKnown("webfinger", m)
}
