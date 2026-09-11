package response

import "runtime"

type SafeError struct {
	Code     Code
	Internal error
	Caller   string

	pcs []uintptr
}

func NewSafeError(code Code, internal error, caller string) *SafeError {

	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	return &SafeError{Code: code, Internal: internal, Caller: caller, pcs: pcs[:n]}
}

func (this *SafeError) Error() string {
	return this.Internal.Error()
}

func (this *SafeError) StackTrace() []uintptr {
	return this.pcs
}
