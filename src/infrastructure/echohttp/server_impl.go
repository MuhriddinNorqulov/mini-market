package echohttp

import (
	"context"
	"errors"
	"fmt"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/infrastructure/echohttp/defaults"
	"mini-market/src/infrastructure/echohttp/mapper"
	"mini-market/src/infrastructure/env"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
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
	return err
}

func (this *EchoServerImpl) Shutdown(ctx context.Context) error {
	return this.echo.Shutdown(ctx)
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
