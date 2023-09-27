package routes

const (
	UserParameterName = "id"
)

type Actor struct {
	root string
}

var _ Partial = Actor{}
var _ Parameterized = Actor{}

func (u Actor) PartialRoute() string {
	return u.root + "/:" + UserParameterName
}

func (u Actor) FullRoute(id string) string {
	return u.root + "/" + id
}
