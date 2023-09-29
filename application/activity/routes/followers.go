package routes

type Followers struct {
	root string
}

func (r Followers) Route() Route {
	return Route{root: r.root, routeName: "followers"}
}
