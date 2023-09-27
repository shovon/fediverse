package routes

const (
	ActorsRoute = "actors"
)

type Actors struct {
	root string
}

var _ Route = Actors{}

func (u Actors) FullRoute() string {
	return u.root + "/" + ActorsRoute
}

func (u Actors) Actor() Actor {
	return Actor{u.FullRoute()}
}
