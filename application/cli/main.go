package main

import (
	"fediverse/application/crypto"
	"fediverse/application/post"
	"flag"
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
	case "genrsa":
		pair, err := crypto.GenerateRSAKeyPair(2048)
		if err != nil {
			panic(err)
		}

		fs := flag.NewFlagSet("genrsa", flag.ExitOnError)
		var showPublic bool
		fs.BoolVar(&showPublic, "public", false, "Show public key")
		fs.Parse(os.Args[2:])

		pemPair := crypto.ToPemPair(pair)
		fmt.Println(string(pemPair.PrivateKey))
		if showPublic {
			fmt.Println(string(pemPair.PublicKey))
		}
	default:
		fmt.Println("Unknown command " + command + ". Expecting a `create` command")
	}
}
