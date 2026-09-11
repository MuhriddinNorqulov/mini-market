package httpport

import (
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type HandlerFunc func(c ctx.Context) error

type IHandler interface {
	Handle(c ctx.Context) error
}
