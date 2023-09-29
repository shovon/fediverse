package routes

type Actors struct {
	root string
}

func (r Actors) Route() Route {
	return Route{root: r.root, routeName: "actors"}
}

func (r Actors) Actor() Actor {
	return Actor{root: r.Route().FullRoute()}
}
