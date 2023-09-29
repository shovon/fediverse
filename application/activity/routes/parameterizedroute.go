package routes

type ParameterizedRouter struct {
	root          string
	parameterName string
}

func (r ParameterizedRouter) FullRoute(value string) string {
	return r.root + "/" + value
}

func (r ParameterizedRouter) ParameterizedRoute() string {
	return r.root + "/:" + r.parameterName
}

func (r ParameterizedRouter) ParameterName() string {
	return r.parameterName
}

type ParameterizedRouteGetter[T RouteGetter] struct {
	parameterizedRouter ParameterizedRouter
	fullRoute           func(root string) T
}

func (r ParameterizedRouteGetter[T]) FullRoute(value string) T {
	return r.fullRoute(r.parameterizedRouter.FullRoute(value))
}

func (r ParameterizedRouteGetter[T]) ParameterizedRoute() T {
	return r.fullRoute(r.parameterizedRouter.ParameterizedRoute())
}
