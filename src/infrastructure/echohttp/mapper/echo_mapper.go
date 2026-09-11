package mapper

import (
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/infrastructure/echohttp/context"

	"github.com/labstack/echo/v4"
)

func ToEchoHandler(h httpport.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(mustLookup(c))
	}
}

func ToHTTPHandler(h echo.HandlerFunc) httpport.HandlerFunc {
	return func(c ctx.Context) error {

		return h(c.(*context.EchoContext).Unwrap())
	}
}

func ToEchoMiddlewareMapper(middleware httpport.Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return middleware(ToHTTPHandler(next))(mustLookup(c))
		}
	}
}

func mustLookup(c echo.Context) ctx.Context {
	wrapped, ok := context.Lookup(c)
	if !ok {
		panic("echohttp/mapper: no *context.EchoContext on request; is context.ContextMiddleware registered?")
	}
	return wrapped
}

func MiddlewareListMapper(middlewares []httpport.Middleware) []echo.MiddlewareFunc {
	mws := make([]echo.MiddlewareFunc, len(middlewares))
	for i, m := range middlewares {
		mws[i] = ToEchoMiddlewareMapper(m)
	}
	return mws
}
