package main

import (
	"context"
	"flag"
	"log"

	container "mini-market/cmd/container"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/entrypoint/seed"
	"mini-market/src/infrastructure/env"
)

var mode string

func init() {
	flag.StringVar(&mode, "mode", "http", "run mode: http | async | check-env | seed-users")
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
	case "seed-users":
		e := env.NewEnv()
		app := container.InitSeedApp()
		specs := []seed.DefaultUserSpec{
			{Username: e.DefaultAdminUsername, Password: e.DefaultAdminPassword, Role: enum.RoleAdmin},
			{Username: e.DefaultUserUsername, Password: e.DefaultUserPassword, Role: enum.RoleUser},
			{Username: e.DefaultDeveloperUsername, Password: e.DefaultDeveloperPassword, Role: enum.RoleDeveloper},
		}
		if err := app.SeedDefaultUsers(context.Background(), specs); err != nil {
			log.Fatalf("seed default users: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %q", mode)
	}
}
