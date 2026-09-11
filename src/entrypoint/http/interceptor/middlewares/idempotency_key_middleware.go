package middlewares

import (
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

func IdempotencyKeyMiddleware(next httpport.HandlerFunc) httpport.HandlerFunc {
	return func(c ctx.Context) error {
		idempotencyKey := c.GetRequest().Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			return response.NewResponse(response.CodeBadRequest, false, nil, "Idempotency-Key header is required")
		}

		c.SetIdempotencyKey(idempotencyKey)
		return next(c)
	}
}
