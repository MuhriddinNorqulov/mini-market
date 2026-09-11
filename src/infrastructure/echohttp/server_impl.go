package echohttp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/infrastructure/echohttp/defaults"
	"mini-market/src/infrastructure/echohttp/mapper"
	"mini-market/src/infrastructure/env"
	"mini-market/src/infrastructure/logger"
	"mini-market/src/infrastructure/sentry"
	"mini-market/src/infrastructure/telemetry"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"
)

// @inject
func NewEcho(validator *defaults.RequestValidator) *echo.Echo {
	e := echo.New()
	e.Validator = validator
	return e
}

type EchoServerImpl struct {
	echo                         *echo.Echo
	env                          *env.Env
	developmentGroup             *defaults.DevelopmentGroup
	developerBasicAuthMiddleware *defaults.DeveloperUserBasicAuthMiddleware
	recoveryMiddleware           *defaults.RecoveryMiddleware
	httpLoggerMiddleware         *defaults.HttpLoggerMiddleware
	consoleLoggerMiddleware      *defaults.ConsoleLoggerMiddleware
	errorRecorderMiddleware      *defaults.ErrorRecorderMiddleware
	logger                       *logger.HttpLogger
	telemetry                    *telemetry.Telemetry
	sentry                       *sentry.Client
}

// @inject
func NewEchoServerImpl(
	echo *echo.Echo,
	cfg *env.Env,
	developmentGroup *defaults.DevelopmentGroup,
	developerBasicAuthMiddleware *defaults.DeveloperUserBasicAuthMiddleware,
	recoveryMiddleware *defaults.RecoveryMiddleware,
	httpLoggerMiddleware *defaults.HttpLoggerMiddleware,
	consoleLoggerMiddleware *defaults.ConsoleLoggerMiddleware,
	errorRecorderMiddleware *defaults.ErrorRecorderMiddleware,
	log *logger.HttpLogger,
	tel *telemetry.Telemetry,
	sentryClient *sentry.Client,
) httpport.HTTPServer {
	return &EchoServerImpl{
		echo:                         echo,
		env:                          cfg,
		developmentGroup:             developmentGroup,
		developerBasicAuthMiddleware: developerBasicAuthMiddleware,
		recoveryMiddleware:           recoveryMiddleware,
		httpLoggerMiddleware:         httpLoggerMiddleware,
		consoleLoggerMiddleware:      consoleLoggerMiddleware,
		errorRecorderMiddleware:      errorRecorderMiddleware,
		logger:                       log,
		telemetry:                    tel,
		sentry:                       sentryClient,
	}
}

func (this *EchoServerImpl) Use(middlewares ...httpport.Middleware) {
	_mws := mapper.MiddlewareListMapper(middlewares)
	this.echo.Use(_mws...)
}

func (this *EchoServerImpl) Run() error {
	err := this.echo.Start(fmt.Sprintf(":%s", this.env.HttpPort))
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	if err != nil {
		this.logger.Error("http server error", zap.Error(err))
	}
	return err
}

func (this *EchoServerImpl) Shutdown(ctx context.Context) error {
	err := this.echo.Shutdown(ctx)

	telemetryCtx, telemetryCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer telemetryCancel()
	if tErr := this.telemetry.Shutdown(telemetryCtx); tErr != nil {
		log.Printf("[telemetry] shutdown: %v", tErr)
	}

	sentryCtx, sentryCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer sentryCancel()
	if sErr := this.sentry.Shutdown(sentryCtx); sErr != nil {
		log.Printf("[sentry] shutdown: %v", sErr)
	}

	return err
}

func (this *EchoServerImpl) Group(prefix string, middlewares ...httpport.Middleware) httpport.Group {
	_mws := mapper.MiddlewareListMapper(middlewares)
	group := this.echo.Group(prefix, _mws...)
	return NewEchoGroupImpl(group)
}

func (this *EchoServerImpl) Init() {

	this.echo.Use(otelecho.Middleware(this.env.OtelServiceNameHTTP))

	this.echo.Use(this.recoveryMiddleware.Wrap)
	this.echo.Use(this.consoleLoggerMiddleware.Wrap)

	this.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: this.env.CorsAllowOrigin,

		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "Idempotency-Key", "traceparent", "tracestate", "baggage"},
		AllowCredentials: this.env.CorsAllowCredentials,
	}))

	this.echo.Use(this.httpLoggerMiddleware.Wrap)
	this.echo.Use(this.errorRecorderMiddleware.Wrap)

	this.echo.Use(defaults.ContextMiddleware)

	this.developmentGroup.RegisterRoutes(this.echo.Group("/dev", this.developerBasicAuthMiddleware.Wrap))
}
