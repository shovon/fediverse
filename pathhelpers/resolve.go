package pathhelpers

import (
	"path"
)

func Resolve(paths ...string) string {
	for i := len(paths) - 1; i >= 0; i-- {
		p := paths[i]
		if len(p) == 0 {
			continue
		}
		if p[0] == '/' {
			paths = paths[i:]
			break
		}
	}
	return path.Join(paths...)
}
