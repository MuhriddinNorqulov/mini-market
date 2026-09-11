package groups

import (
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/entrypoint/http/handlers/auth"
)

type AuthGroup struct {
	registerHandler *auth.RegisterHandler
	loginHandler    *auth.LoginHandler
}

// @inject
func NewAuthGroup(
	registerHandler *auth.RegisterHandler,
	loginHandler *auth.LoginHandler,
) *AuthGroup {
	return &AuthGroup{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
	}
}

func (this *AuthGroup) RegisterRoutes(g httpport.Group) {
	g.POST("/register", this.registerHandler.Handle)
	g.POST("/login", this.loginHandler.Handle)
}
