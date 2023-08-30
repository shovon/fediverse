package rfc3230

type Digester interface {
	// Token returns a constant that represents the token or digest name, such as
	// sha-256 or sha-512.
	Token() string

	// Digest returns the digest of the body as a string.
	Digest([]byte) (string, error)
}

type digester struct {
	token  string
	digest func([]byte) (string, error)
}

var _ Digester = digester{}

func (d digester) Token() string {
	return d.token
}

func (d digester) Digest(body []byte) (string, error) {
	return d.digest(body)
}

func CreateDigester(token string, digest func([]byte) (string, error)) Digester {
	return digester{
		token:  token,
		digest: digest,
	}
}
