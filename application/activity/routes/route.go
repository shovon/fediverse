package routes

type Route interface {
	FullRoute() string
}

type Parameterized interface {
	ParameterName() string
	RouteSubbed(string) string
}
