package http

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/entrypoint/http/interceptor/middlewares"
	"mini-market/src/infrastructure/env"
	"mini-market/src/infrastructure/logger"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"

	"go.uber.org/zap"
)

type App struct {
	server              httpport.HTTPServer
	env                 *env.Env
	jwtAuthMiddleware   *middlewares.JwtAuthMiddleware
	responseMiddleware  *middlewares.ResponseMiddleware
	rateLimitMiddleware *middlewares.RateLimitMiddleware
	telemetry           *telemetry.Telemetry
	sentry              *sentry.Client
	logger              *logger.HttpLogger
}

// @inject
func NewApp(
	server httpport.HTTPServer,
	env *env.Env,
	jwtAuthMiddleware *middlewares.JwtAuthMiddleware,
	responseMiddleware *middlewares.ResponseMiddleware,
	rateLimitMiddleware *middlewares.RateLimitMiddleware,
	telemetry *telemetry.Telemetry,
	sentry *sentry.Client,
	log *logger.HttpLogger,
) *App {
	return &App{
		server:              server,
		env:                 env,
		jwtAuthMiddleware:   jwtAuthMiddleware,
		responseMiddleware:  responseMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		telemetry:           telemetry,
		sentry:              sentry,
		logger:              log,
	}
}

func (this *App) Init() {
	this.server.Init()
	this.server.Use(this.jwtAuthMiddleware.Wrap)
	this.server.Use(this.responseMiddleware.Wrap)
}

func (this *App) Start() {
	errCh := make(chan error, 1)
	go func() { errCh <- this.server.Run() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var runErr error
	select {
	case err := <-errCh:
		if err != nil {
			runErr = err
			this.logger.Error("http server error", zap.Error(err))
		}
	case <-stop:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = this.server.Shutdown(shutdownCtx)

	telemetryCtx, telemetryCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer telemetryCancel()
	if err := this.telemetry.Shutdown(telemetryCtx); err != nil {
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
