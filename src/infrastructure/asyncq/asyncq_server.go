package async

import (
	"context"
	"log"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/async"
	"mini-market/src/infrastructure/asyncq/mapper"
	"mini-market/src/infrastructure/asyncq/middlewares"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"
	"time"

	"github.com/hibiken/asynq"
)

type AsynqServerImpl struct {
	s   *asynq.Server
	mux *asynq.ServeMux

	logMiddleware   *middlewares.TaskLogMiddleware
	alertMiddleware *middlewares.TaskAlertMiddleware
	traceMiddleware *middlewares.TaskTraceMiddleware

	telemetry *telemetry.Telemetry
	sentry    *sentry.Client
}

// @inject
func NewAsynqServerImpl(
	s *asynq.Server,
	mux *asynq.ServeMux,
	logMiddleware *middlewares.TaskLogMiddleware,
	alertMiddleware *middlewares.TaskAlertMiddleware,
	traceMiddleware *middlewares.TaskTraceMiddleware,
	tel *telemetry.Telemetry,
	sentryClient *sentry.Client,
) async.AsyncServer {
	return &AsynqServerImpl{
		s:               s,
		mux:             mux,
		logMiddleware:   logMiddleware,
		alertMiddleware: alertMiddleware,
		traceMiddleware: traceMiddleware,
		telemetry:       tel,
		sentry:          sentryClient,
	}
}

func (this *AsynqServerImpl) Init() {

	this.Use(this.traceMiddleware.Wrap)
	this.Use(this.alertMiddleware.Wrap)
	this.Use(this.logMiddleware.Wrap)
}

func (this *AsynqServerImpl) Run() error {
	return this.s.Run(this.mux)
}

func (this *AsynqServerImpl) Shutdown(ctx context.Context) error {
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

	return nil
}

func (this *AsynqServerImpl) Use(middlewares ...async.AsyncTaskMiddleware) {
	for _, middleware := range middlewares {
		this.mux.Use(mapper.AsynqTaskMiddlewareMapper(middleware))
	}
}

func (this *AsynqServerImpl) HandlerFunc(taskName enum.TaskType, handler async.AsyncTaskHandler) {
	this.mux.HandleFunc(string(taskName), mapper.ToAsynQTask(handler).ProcessTask)
}

func (this *AsynqServerImpl) ProcessTask(ctx context.Context, task *async.Task) error {
	return this.mux.ProcessTask(ctx, asynq.NewTask(string(task.TaskType), task.Payload))
}
