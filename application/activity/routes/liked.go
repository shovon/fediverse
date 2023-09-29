package routes

type Liked struct {
	root string
}

func (r Liked) Route() Route {
	return Route{root: r.root, routeName: "inbox"}
}
