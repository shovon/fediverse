package jsonldcontext

type IRIContext string

var _ ValidContext = IRIContext("")

func (i IRIContext) uselessValidContext() useless {
	return useless{}
}
