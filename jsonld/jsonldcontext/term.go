package jsonldcontext

type Term interface {
	uselessTerm() useless
	MarshalJSON() ([]byte, error)
}
