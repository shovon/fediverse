package webfinger

import (
	"fediverse/jrd"
	"fediverse/nullable"
	"net/http"
)

func HandleRel(j jrd.JRD, r *http.Request) jrd.JRD {
	rel := r.URL.Query().Get("rel")
	if rel != "" {
		currentLinks, ok := j.Links.Value()
		if ok {
			links := []jrd.Link{}
			for _, link := range currentLinks {
				if link.Rel == rel {
					links = append(links, link)
				}
			}
			j.Links = nullable.Just(links)
		}
	}
	return j
}
