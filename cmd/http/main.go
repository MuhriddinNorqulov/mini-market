package main

import container "mini-market/cmd/container"

// @title Mini Market API
// @version 1.0
// @BasePath /api/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.basic BasicAuth
func main() {
	app := container.InitHttpApp()
	app.Init()
	app.Start()
}
