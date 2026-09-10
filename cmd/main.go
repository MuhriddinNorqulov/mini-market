package main

import (
	"flag"
	"log"

	container "mini-market/cmd/container"
	"mini-market/src/infrastructure/env"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "http", "run mode: http | async | check-env")
}

func main() {
	flag.Parse()
	switch mode {
	case "http":
		app := container.InitHttpApp()
		app.Init()
		app.Start()
	case "async":
		app := container.InitAsyncApp()
		app.Init()
		app.Start()
	case "check-env":
		if err := env.Check(); err != nil {
			log.Fatalf("env check: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %q", mode)
	}
}
