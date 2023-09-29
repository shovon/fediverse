package routes

type SharedInbox struct {
	root string
}

func (r SharedInbox) Route() Route {
	return Route{root: r.root, routeName: "sharedInbox"}
}
