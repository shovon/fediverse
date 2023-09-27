package routes

const (
	InboxRoute = "inbox"
)

type Inbox struct {
	root string
}

var _ Route = Inbox{}

func (i Inbox) FullRoute() string {
	return i.root + "/" + InboxRoute
}
