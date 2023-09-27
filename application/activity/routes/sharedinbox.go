package routes

const SharedInboxRoute = "sharedinbox"

type SharedInbox struct {
	root string
}

var _ Route = SharedInbox{}

func (s SharedInbox) FullRoute() string {
	return s.root + "/" + SharedInboxRoute
}

func (s SharedInbox) PartialRoute() string {
	return "/" + SharedInboxRoute
}
