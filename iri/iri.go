package iri

import (
	"net/url"
	"strings"
)

// TODO: unit test this
type IRI struct {
	Scheme      string
	Opaque      string        // encoded opaque data
	User        *url.Userinfo // username and password information
	Host        string        // host or host:port
	Path        string        // path (relative paths may omit leading slash)
	RawPath     string        // encoded path hint (see EscapedPath method)
	OmitHost    bool          // do not emit empty host (authority)
	ForceQuery  bool          // append a query ('?') even if RawQuery is empty
	RawQuery    string        // encoded query values, without '?'
	Fragment    string        // fragment for references, without '#'
	RawFragment string        // encoded fragment hint (see EscapedFragment method)
}

func ParseIRI(raw string) (IRI, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return IRI{}, err
	}
	return IRI{
		Scheme:      u.Scheme,
		Opaque:      u.Opaque,
		User:        u.User,
		Host:        u.Host,
		Path:        u.Path,
		RawPath:     u.RawPath,
		OmitHost:    u.OmitHost,
		ForceQuery:  u.ForceQuery,
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
	}, nil
}

func (u IRI) String() string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Opaque != "" {
		buf.WriteString(u.Opaque)
	} else {
		if u.Scheme != "" || u.Host != "" || u.User != nil {
			if u.OmitHost && u.Host == "" && u.User == nil {
				// omit empty host
			} else {
				if u.Host != "" || u.Path != "" || u.User != nil {
					buf.WriteString("//")
				}
				if ui := u.User; ui != nil {
					buf.WriteString(ui.String())
					buf.WriteByte('@')
				}
				if h := u.Host; h != "" {
					buf.WriteString(h)
				}
			}
		}
		path := u.RawPath
		if path != "" && path[0] != '/' && u.Host != "" {
			buf.WriteByte('/')
		}
		if buf.Len() == 0 {
			// RFC 3986 §4.2
			// A path segment that contains a colon character (e.g., "this:that")
			// cannot be used as the first segment of a relative-path reference, as
			// it would be mistaken for a scheme name. Such a segment must be
			// preceded by a dot-segment (e.g., "./this:that") to make a relative-
			// path reference.
			if segment, _, _ := strings.Cut(path, "/"); strings.Contains(segment, ":") {
				buf.WriteString("./")
			}
		}
		buf.WriteString(path)
	}
	if u.ForceQuery || u.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(u.RawQuery)
	}
	if u.RawFragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.RawFragment)
	}
	return buf.String()
}
