package httpport

import "context"

type HTTPServer interface {
	Init()
	Run() error
	Shutdown(ctx context.Context) error

	Use(middlewares ...Middleware)
	Group(prefix string, middlewares ...Middleware) Group
}
