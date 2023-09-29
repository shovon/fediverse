package routes

type RouteGetter interface {
	Route() Route
}

type Route struct {
	root      string
	routeName string
}

func (r Route) FullRoute() string {
	return r.root + "/" + r.routeName
}
