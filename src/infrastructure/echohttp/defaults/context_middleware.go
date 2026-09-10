package defaults

import (
	"mini-market/src/infrastructure/echohttp/context"

	"github.com/labstack/echo/v4"
)

func ContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		wrapped := context.NewEchoContext(c)

		wrapped.(*context.EchoContext).Store()
		return next(c)
	}
}
