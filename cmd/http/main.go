package main

import container "mini-market/cmd/container"

// @title Mini Market API
// @version 1.0
// @BasePath /api/
func main() {
	app := container.InitHttpApp()
	app.Init()
	app.Start()
}
