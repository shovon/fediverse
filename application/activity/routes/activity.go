package routes

const (
	ActivityRoute = "activity"
)

type Activity struct {
	root string
}

func (a Activity) FullRoute() string {
	return a.root + ActivityRoute
}

func (a Activity) PartialRoute() string {
	return "/" + ActivityRoute
}

func (a Activity) SharedInbox() SharedInbox {
	return SharedInbox{a.FullRoute()}
}

func (a Activity) Actors() Actors {
	return Actors{a.FullRoute()}
}
