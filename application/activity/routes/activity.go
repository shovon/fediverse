package routes

const (
	ActivityRoute = "activity"
)

type Activity struct {
	root string
}

func (a Activity) FullRoute() string {
	return "/" + ActivityRoute + a.root
}

func (a Activity) PartialRoute() string {
	return "/" + ActivityRoute
}

func (a Activity) SharedInbox() SharedInbox {
	return SharedInbox{a.FullRoute()}
}
