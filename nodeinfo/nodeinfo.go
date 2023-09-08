package nodeinfo

import (
	"fediverse/httphelpers"
	"fediverse/httphelpers/httperrors"
	"fediverse/jrd"
	"fediverse/jrd/jrdhttp"
	"fediverse/nodeinfo/schema2p0"
	"fediverse/nullable"
	"fediverse/wellknown"
	"net/http"
)

type SoftwareInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository string `json:"repository,omitempty"`
	Homepage   string `json:"homepage,omitempty"`
}

type UsersStats struct {
	Total          uint `json:"total"`
	ActiveHalfyear uint `json:"activeHalfyear"`
	ActiveMonth    uint `json:"activeMonth"`
}

type Usage struct {
	Users         UsersStats `json:"users"`
	LocalPosts    uint       `json:"localPosts"`
	LocalComments uint       `json:"localComments"`
}

type NodeInfoProps struct {
	Software          SoftwareInfo `json:"software"`
	OpenRegistrations bool         `json:"openRegistrations"`
	Usage             Usage        `json:"usage"`
}

// CreateNodeInfoMiddleware creates an opinionated middleware for processing
// NodeInfo requests.
//
// The middleware expects a string as a host (such as "https://example.com"),
// the root route for all NodeInfo endpoints (such as "/nodeinfo"), and a
// function that will return some NodeInfo info.
//
// This function is considered "opinionated" because the NodeInfo spec leaves a
// lot to the imagination (such as which routes should—say—provicde the NodeInfo
// v2.0 content be served from), but this function just jams a single idea down
// everyone's throat, by accepting just one root route, and serving all the
// NodeInfo content from there.
//
// If you want more flexibility, this library is not for you.
//
// That said, NodeInfo's purpose is quite simple: tell other agents what a
// particular instance is all about. There shouldn't be much else to it,
// honestly.
//
// But it is understandable that in rare occasions, that there may be situations
// where this library may be a bit too restrictive.
func CreateNodeInfoMiddleware(origin string, nodeInfoRoot string, handler func() NodeInfoProps) func(http.Handler) http.Handler {
	wellKnown := wellknown.WellKnown("nodeinfo", jrdhttp.CreateJRDHandler(func(r *http.Request) (jrd.JRD, httperrors.HTTPError) {
		return jrd.JRD{
			Links: nullable.Just([]jrd.Link{
				{
					Rel:  "http://nodeinfo.diaspora.software/ns/schema/2.0",
					Href: origin + nodeInfoRoot + "/2.0",
				},
			}),
		}, nil
	}))

	schema2p0 := httphelpers.Group(
		nodeInfoRoot, httphelpers.Route("/2.0").Process(httphelpers.ToMiddleware(httphelpers.ToHandlerFunc(httphelpers.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
			nodeInfoProps := handler()
			schema := schema2p0.Schema{
				Software: schema2p0.Software{
					Name:    nodeInfoProps.Software.Name,
					Version: nodeInfoProps.Software.Version,
				},
				Usage: schema2p0.Usage{
					Users: schema2p0.UsersStats{
						Total:          nodeInfoProps.Usage.Users.Total,
						ActiveHalfyear: nodeInfoProps.Usage.Users.ActiveHalfyear,
						ActiveMonth:    nodeInfoProps.Usage.Users.ActiveMonth,
					},
					LocalPosts:    nodeInfoProps.Usage.LocalPosts,
					LocalComments: nodeInfoProps.Usage.LocalComments,
				},
			}
			return httphelpers.WriteJSON(w, schema)
		})))))

	return func(next http.Handler) http.Handler {
		return wellKnown(schema2p0(next))
	}
}
