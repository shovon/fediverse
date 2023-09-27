package routes

type Partial interface {
	PartialRoute() string
}

type Full interface {
	FullRoute() string
}

type Parameterized interface {
	FullRoute(string) string
}
