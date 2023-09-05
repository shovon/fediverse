package main

import (
	"fediverse/application"
	"fediverse/application/schema"
)

func main() {
	err := schema.Initialize()
	if err != nil {
		panic(err)
	}
	panic(application.Start())
}
