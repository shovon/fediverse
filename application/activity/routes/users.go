package routes

const (
	UsersRoute = "users"
)

type Users struct {
	root string
}

var _ Partial = Users{}
var _ Full = Users{}

func (u Users) PartialRoute() string {
	return u.root + "/:" + UserParameterName
}

func (u Users) FullRoute() string {
	return u.root + "/" + UsersRoute
}
