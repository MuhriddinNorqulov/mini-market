package middlewares

import (
	"context"
	"mini-market/src/core/application/services"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/ports/reqctx"
	"mini-market/src/core/domain/ports/security"
	"mini-market/src/infrastructure/telemetry"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type JwtAuthMiddleware struct {
	authService *services.UserAuthTokenService
	denylist    security.SessionDenylist
}

// @inject
func NewJwtAuthMiddleware(authService *services.UserAuthTokenService, denylist security.SessionDenylist) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{authService: authService, denylist: denylist}
}

func (this *JwtAuthMiddleware) Wrap(next httpport.HandlerFunc) httpport.HandlerFunc {
	return func(c ctx.Context) error {
		authHeader := c.GetRequest().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return next(c)
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, sessionID, err := this.authService.VerifyToken(token, enum.TokenTypeAccess)
		if err != nil {
			return next(c)
		}

		if this.isRevoked(c.GetContext(), sessionID) {

			return next(c)
		}

		c.SetUser(user)

		trace.SpanFromContext(c.GetContext()).SetAttributes(
			attribute.Int(telemetry.AttrUserID, int(user.ID)),
			attribute.String(telemetry.AttrUserRole, string(user.Role)),
		)

		c.SetContext(reqctx.WithSessionID(
			reqctx.WithUser(c.GetContext(), user.ID, string(user.Role)),
			sessionID,
		))

		return next(c)
	}
}

func (this *JwtAuthMiddleware) isRevoked(c context.Context, sessionID string) bool {

	if sessionID == "" {
		return false
	}

	revoked, err := this.denylist.IsRevoked(c, sessionID)
	if err != nil {
		trace.SpanFromContext(c).RecordError(err)
		return false
	}
	return revoked
}
