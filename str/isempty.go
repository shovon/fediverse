package str

func IsEmpty(s string) bool {
	return len(TrimWhitespace(s)) == 0
}
