package echohttp

import (
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/infrastructure/echohttp/mapper"

	"github.com/labstack/echo/v4"
)

type EchoGroupImpl struct {
	g *echo.Group
}

func NewEchoGroupImpl(g *echo.Group) httpport.Group {
	return &EchoGroupImpl{g: g}
}

func (e *EchoGroupImpl) Group(prefix string, middlewares ...httpport.Middleware) httpport.Group {
	mws := mapper.MiddlewareListMapper(middlewares)
	return NewEchoGroupImpl(e.g.Group(prefix, mws...))
}

func (e *EchoGroupImpl) Use(middlewares ...httpport.Middleware) {
	mws := mapper.MiddlewareListMapper(middlewares)
	e.g.Use(mws...)
}

func (e *EchoGroupImpl) GET(path string, handler httpport.HandlerFunc, middlewares ...httpport.Middleware) {
	e.g.GET(path, mapper.ToEchoHandler(handler), mapper.MiddlewareListMapper(middlewares)...)
}

func (e *EchoGroupImpl) POST(path string, handler httpport.HandlerFunc, middlewares ...httpport.Middleware) {
	e.g.POST(path, mapper.ToEchoHandler(handler), mapper.MiddlewareListMapper(middlewares)...)
}

func (e *EchoGroupImpl) PUT(path string, handler httpport.HandlerFunc, middlewares ...httpport.Middleware) {
	e.g.PUT(path, mapper.ToEchoHandler(handler), mapper.MiddlewareListMapper(middlewares)...)
}

func (e *EchoGroupImpl) DELETE(path string, handler httpport.HandlerFunc, middlewares ...httpport.Middleware) {
	e.g.DELETE(path, mapper.ToEchoHandler(handler), mapper.MiddlewareListMapper(middlewares)...)
}

func (e *EchoGroupImpl) PATCH(path string, handler httpport.HandlerFunc, middlewares ...httpport.Middleware) {
	e.g.PATCH(path, mapper.ToEchoHandler(handler), mapper.MiddlewareListMapper(middlewares)...)
}

func (e *EchoGroupImpl) Any(path string, handler httpport.HandlerFunc, middlewares ...httpport.Middleware) {
	e.g.Any(path, mapper.ToEchoHandler(handler), mapper.MiddlewareListMapper(middlewares)...)
}
