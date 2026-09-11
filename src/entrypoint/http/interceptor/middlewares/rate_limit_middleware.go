package middlewares

import (
	"fmt"
	"time"

	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/ports/ratelimit"

	"go.opentelemetry.io/otel/trace"
)

type RateLimitMiddleware struct {
	limiter ratelimit.Limiter
}

// @inject
func NewRateLimitMiddleware(limiter ratelimit.Limiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: limiter}
}

func (this *RateLimitMiddleware) PerUser(name string, limit int, window time.Duration) httpport.Middleware {
	return func(next httpport.HandlerFunc) httpport.HandlerFunc {
		return func(c ctx.Context) error {
			user := c.User()
			if user == nil {
				return next(c)
			}

			key := fmt.Sprintf("ratelimit:%s:user:%d", name, user.ID)

			allowed, err := this.limiter.Allow(c.GetContext(), key, limit, window)
			if err != nil {
				trace.SpanFromContext(c.GetContext()).RecordError(err)
				return next(c)
			}

			if !allowed {
				return response.NewResponse(
					response.CodeTooManyRequests, false, nil,
					"Juda ko'p so'rov yuborildi, biroz kutib qayta urinib ko'ring",
				)
			}

			return next(c)
		}
	}
}
