package routes

type Actor struct {
	root string
}

func (r Actor) Route() ParameterizedRouter {
	return ParameterizedRouter{root: r.root, parameterName: "id"}
}

func (r Actor) Inbox() ParameterizedRouteGetter[Inbox] {
	return ParameterizedRouteGetter[Inbox]{
		parameterizedRouter: r.Route(),
		fullRoute: func(root string) Inbox {
			return Inbox{root: root}
		},
	}
}

func (r Actor) Outbox() ParameterizedRouteGetter[Outbox] {
	return ParameterizedRouteGetter[Outbox]{
		parameterizedRouter: r.Route(),
		fullRoute: func(root string) Outbox {
			return Outbox{root: root}
		},
	}
}

func (r Actor) Following() ParameterizedRouteGetter[Following] {
	return ParameterizedRouteGetter[Following]{
		parameterizedRouter: r.Route(),
		fullRoute: func(root string) Following {
			return Following{root: root}
		},
	}
}

func (r Actor) Followers() ParameterizedRouteGetter[Followers] {
	return ParameterizedRouteGetter[Followers]{
		parameterizedRouter: r.Route(),
		fullRoute: func(root string) Followers {
			return Followers{root: root}
		},
	}
}
