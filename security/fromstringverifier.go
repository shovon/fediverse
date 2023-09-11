package security

// FromStringVerifier is a type that allows you to verify a signature in the
// form of a string.
//
// This exists precisely because in certain protocols, the signature absolutely
// must be encoded into a UTF-8 string, and not a byte slice.
//
// What container format to use varies from application to application. Hence,
// why it's not a good idea to excusively return a byte slice, but instead to
// return a string.
type FromStringVerifier interface {
	Verify(payload []byte, signature string) error
}
