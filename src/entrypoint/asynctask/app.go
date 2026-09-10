package asynctask

import (
	"context"
	"log"
	"time"

	"mini-market/src/core/domain/ports/async"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"
)

type App struct {
	server    async.AsyncServer
	telemetry *telemetry.Telemetry
	sentry    *sentry.Client
}

// @inject
func NewAsyncApp(
	server async.AsyncServer,
	telemetry *telemetry.Telemetry,
	sentry *sentry.Client,
) *App {
	return &App{server: server, telemetry: telemetry, sentry: sentry}
}

func (this *App) Init() {
	this.server.Init()
}

func (this *App) Start() {
	runErr := this.server.Run()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := this.telemetry.Shutdown(shutdownCtx); err != nil {
		log.Printf("[telemetry] shutdown: %v", err)
	}

	sentryCtx, sentryCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer sentryCancel()
	if err := this.sentry.Shutdown(sentryCtx); err != nil {
		log.Printf("[sentry] shutdown: %v", err)
	}

	if runErr != nil {
		panic(runErr)
	}
}
