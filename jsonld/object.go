package jsonld

type useless struct{}

type ContextType[T any] interface {
	useless() useless
	Value() T
}

type Slice struct {
	Contexts []ContextType[any] `json:"@context"`
}

var _ ContextType[[]ContextType[any]] = Slice{}

func (s Slice) useless() useless {
	return useless{}
}

func (s Slice) Value() []ContextType[any] {
	return s.Contexts
}

type Object struct {
	Context ContextType[any] `json:"@context"`
	Data    any
}
