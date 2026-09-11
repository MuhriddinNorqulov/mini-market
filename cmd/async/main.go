package main

import container "mini-market/cmd/container"

func main() {
	app := container.InitAsyncApp()
	app.Init()
	app.Start()
}
