package routes

type Inbox struct {
	root string
}

func (r Inbox) Route() Route {
	return Route{root: r.root, routeName: "inbox"}
}
