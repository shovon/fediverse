package pathhelpers

import (
	"strings"
)

func Match(pattern, path string) (bool, map[string]string) {
	patternSplit := strings.Split(pattern, "/")
	pathSplit := strings.Split(path, "/")

	if len(patternSplit) != len(pathSplit) {
		return false, nil
	}

	params := make(map[string]string)

	for i, patternPart := range patternSplit {
		pathPart := pathSplit[i]
		// TODO: also add regex checks.
		if strings.HasPrefix(patternPart, ":") {
			params[patternPart[1:]] = pathPart
			continue
		}
		if pathPart != patternPart {
			return false, nil
		}
	}

	return true, params
}

// TODO: unit test this

func PartialMatch(pattern, path string) (bool, string, map[string]string) {
	patternSplit := strings.Split(pattern, "/")
	pathSplit := strings.Split(path, "/")

	if len(patternSplit) > len(pathSplit) {
		return false, "", nil
	}

	params := make(map[string]string)

	for i, patternPart := range patternSplit {
		pathPart := pathSplit[i]
		if strings.HasPrefix(patternPart, ":") {
			params[patternPart[1:]] = pathPart
			continue
		}
		if pathPart != patternPart {
			return false, "", nil
		}
	}

	return true, "/" + strings.Join(pathSplit[len(patternSplit):], "/"), params
}
