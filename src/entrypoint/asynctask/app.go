package asynctask

import (
	"context"
	"time"

	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/async"
	"mini-market/src/entrypoint/asynctask/handlers"
)

type App struct {
	server                 async.AsyncServer
	orderAutoCancelHandler *handlers.OrderAutoCancelHandler
}

// @inject
func NewAsyncApp(
	server async.AsyncServer,
	orderAutoCancelHandler *handlers.OrderAutoCancelHandler,
) *App {
	return &App{server: server, orderAutoCancelHandler: orderAutoCancelHandler}
}

func (this *App) Init() {
	this.server.Init()
	this.server.HandlerFunc(enum.TaskTypeOrderAutoCancel, this.orderAutoCancelHandler)
}

func (this *App) Start() {
	runErr := this.server.Run()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = this.server.Shutdown(shutdownCtx)

	if runErr != nil {
		panic(runErr)
	}
}
