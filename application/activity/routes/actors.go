package routes

const (
	UsersRoute = "actors"
)

type Actors struct {
	root string
}

var _ Partial = Actors{}
var _ Full = Actors{}

func (u Actors) PartialRoute() string {
	return u.root + "/:" + UsersRoute
}

func (u Actors) FullRoute() string {
	return u.root + "/" + UsersRoute
}

func (u Actors) Actor() Actor {
	return Actor{u.FullRoute()}
}
