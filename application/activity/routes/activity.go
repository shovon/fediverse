package routes

type Activity struct {
	root string
}

func (r Activity) Route() Route {
	return Route{root: r.root, routeName: "activity"}
}

func (r Activity) SharedInbox() SharedInbox {
	return SharedInbox{root: r.Route().FullRoute()}
}

func (r Activity) Actors() Actors {
	return Actors{root: r.Route().FullRoute()}
}
