package main

import (
	"fediverse/application/post"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Please provide a command")
		os.Exit(1)
		return
	}

	command := args[0]

	switch command {
	case "create":
		if len(args) <= 1 {
			fmt.Println("Please provide some content")
			os.Exit(1)
			return
		}
		post.CreatePost(args[1])
		fmt.Println("Post successfully created!")
	default:
		fmt.Println("Unknown command " + command + ". Expecting a `create` command")
	}
}
