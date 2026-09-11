package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/entrypoint/http/groups"
	"mini-market/src/entrypoint/http/interceptor/middlewares"
)

type App struct {
	server              httpport.HTTPServer
	jwtAuthMiddleware   *middlewares.JwtAuthMiddleware
	responseMiddleware  *middlewares.ResponseMiddleware
	rateLimitMiddleware *middlewares.RateLimitMiddleware
	productGroup        *groups.ProductGroup
	orderGroup          *groups.OrderGroup
	authGroup           *groups.AuthGroup
	meGroup             *groups.MeGroup
}

// @inject
func NewApp(
	server httpport.HTTPServer,
	jwtAuthMiddleware *middlewares.JwtAuthMiddleware,
	responseMiddleware *middlewares.ResponseMiddleware,
	rateLimitMiddleware *middlewares.RateLimitMiddleware,
	productGroup *groups.ProductGroup,
	orderGroup *groups.OrderGroup,
	authGroup *groups.AuthGroup,
	meGroup *groups.MeGroup,
) *App {
	return &App{
		server:              server,
		jwtAuthMiddleware:   jwtAuthMiddleware,
		responseMiddleware:  responseMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		productGroup:        productGroup,
		orderGroup:          orderGroup,
		authGroup:           authGroup,
		meGroup:             meGroup,
	}
}

func (this *App) Init() {
	this.server.Init()
	this.initMiddlewares()
	this.initGroups()
}

func (this *App) initMiddlewares() {
	this.server.Use(this.jwtAuthMiddleware.Wrap)
	this.server.Use(this.responseMiddleware.Wrap)
}

func (this *App) initGroups() {
	api := this.group("/api")

	this.productGroup.RegisterRoutes(api.Group("/products"))
	this.orderGroup.RegisterRoutes(api.Group("/orders"))
	this.authGroup.RegisterRoutes(api.Group("/auth"))
	this.meGroup.RegisterRoutes(api.Group("/me"))
}

func (this *App) group(prefix string, mws ...httpport.Middleware) httpport.Group {
	return this.server.Group(prefix, mws...)
}

func (this *App) Start() {
	errCh := make(chan error, 1)
	go func() { errCh <- this.server.Run() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var runErr error
	select {
	case err := <-errCh:
		runErr = err
	case <-stop:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = this.server.Shutdown(shutdownCtx)

	if runErr != nil {
		panic(runErr)
	}
}
