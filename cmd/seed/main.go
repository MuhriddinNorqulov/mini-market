package main

import (
	"fmt"
	"os"

	container "mini-market/cmd/container"
)

func main() {
	username := "seed-user"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	app := container.InitSeedApp()
	token, err := app.IssueTestToken(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(token)
}
