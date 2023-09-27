package routes

type Root struct{}

var _ Full = Root{}

func (r Root) Activity() Activity {
	return Activity{r.FullRoute()}
}

func (r Root) FullRoute() string {
	return "/"
}
