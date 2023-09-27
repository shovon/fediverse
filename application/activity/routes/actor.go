package routes

import "fmt"

const (
	UserParameterName = "id"
)

type Actor struct {
	root string
}

var _ Route = Actor{}
var _ Parameterized = Actor{}

func (u Actor) FullRoute() string {
	return u.root + "/:" + UserParameterName
}

func (u Actor) ParameterName() string {
	return UserParameterName
}

func (u Actor) RouteSubbed(id string) string {
	fmt.Println("u.root", u.root)
	return u.root + "/" + id
}
