package requestimpl

import (
	"mini-market/src/core/domain/ports/httpport/request"
	"net/http"
)

type HeaderImpl struct {
	super http.Header
}

func NewHeaderImpl(super http.Header) request.Header {
	return &HeaderImpl{super: super}
}

func (h *HeaderImpl) Get(key string) string {
	return h.super.Get(key)
}

func (h *HeaderImpl) Set(key string, value string) {
	h.super.Set(key, value)
}

func (h *HeaderImpl) Add(key, value string) {
	h.super.Add(key, value)
}
