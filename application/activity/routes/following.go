package routes

type Following struct {
	root string
}

func (r Following) Route() Route {
	return Route{root: r.root, routeName: "following"}
}
