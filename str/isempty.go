package str

import "strings"

// IsEmpty tests whether a string is empty.
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
