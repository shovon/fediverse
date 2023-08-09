package jsonldcontext

type IRITerm string

var _ Term = IRITerm("")

func (i IRITerm) uselessTerm() useless {
	return useless{}
}

func (i IRITerm) MarshalJSON() ([]byte, error) {
	return []byte(i), nil
}
