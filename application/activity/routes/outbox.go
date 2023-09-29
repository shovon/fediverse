package routes

type Outbox struct {
	root string
}

func (r Outbox) Route() Route {
	return Route{root: r.root, routeName: "outbox"}
}
