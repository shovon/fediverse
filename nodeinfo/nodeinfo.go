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

	schema2p0 := httphelpers.Route(nodeInfoRoot+"/2.0", httphelpers.ToMiddleware(httphelpers.ToHandlerFunc(httphelpers.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
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
		return httphelpers.WriteJSON(w, 200, schema, nullable.Null[string]())
	}))))

	return func(next http.Handler) http.Handler {
		return wellKnown(schema2p0(next))
	}
}
