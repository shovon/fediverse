package namespacing

type Prefix string

func (p Prefix) Suffix(value string) string {
	return string(p) + value
}
