package routes

const (
	UserParameterName = "id"
)

type User struct {
	root string
}

var _ Partial = User{}
var _ Parameterized = User{}

func (u User) PartialRoute() string {
	return u.root + "/:" + UserParameterName
}

func (u User) FullRoute(id string) string {
	return u.root + "/" + id
}
