package routes

const (
	InboxRoute = "inbox"
)

type Inbox struct {
	root string
}

var _ Partial = Inbox{}
var _ Full = Inbox{}

func (i Inbox) FullRoute() string {
	return i.root + "/" + InboxRoute
}

func (i Inbox) PartialRoute() string {
	return "/" + InboxRoute
}
