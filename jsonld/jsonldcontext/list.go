package jsonldcontext

type List []ValidContext

var _ ValidContext = List{}

func (l List) uselessValidContext() useless {
	return useless{}
}
