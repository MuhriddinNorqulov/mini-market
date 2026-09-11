package telemetry

import (
	"context"
	"errors"
	"runtime"

	"mini-market/src/core/application/response"
)

type PanicError struct {
	Err error

	pcs []uintptr
}

func NewPanicError(err error, skip int) *PanicError {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip, pcs)
	return &PanicError{Err: err, pcs: pcs[:n]}
}

func (this *PanicError) Error() string { return this.Err.Error() }
func (this *PanicError) Unwrap() error { return this.Err }

func (this *PanicError) StackTrace() []uintptr { return this.pcs }

func ErrorKindOf(err error) (kind string, failure bool) {
	if err == nil {
		return "", false
	}

	if errors.Is(err, context.Canceled) {
		return "", false
	}
	if se, ok := errors.AsType[*response.SafeError](err); ok && errors.Is(se.Internal, context.Canceled) {
		return "", false
	}

	if _, ok := errors.AsType[*PanicError](err); ok {
		return ErrorKindPanic, true
	}

	if _, ok := errors.AsType[*response.Response](err); ok {
		return ErrorKindBusiness, false
	}

	return ErrorKindError, true
}
