package defaults

import (
	"mini-market/src/infrastructure/env"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DeveloperUserBasicAuthMiddleware struct {
	env *env.Env
}

// @inject
func NewDeveloperUserBasicAuthMiddleware(e *env.Env) *DeveloperUserBasicAuthMiddleware {
	return &DeveloperUserBasicAuthMiddleware{env: e}
}

func (this *DeveloperUserBasicAuthMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		return username == this.env.DeveloperBasicAuthUsername && password == this.env.DeveloperBasicAuthPassword, nil
	})(next)
}
