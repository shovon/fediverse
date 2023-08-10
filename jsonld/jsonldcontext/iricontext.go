package jsonldcontext

type IRIContext string

var _ ValidContext = IRIContext("")

func (i IRIContext) uselessValidContext() useless {
	return useless{}
}

func (i IRIContext) MarshalJSON() ([]byte, error) {
	return []byte(i), nil
}
