package defaults

import (
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/security"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DeveloperUserBasicAuthMiddleware struct {
	userRepo repository.UserRepository
	hasher   security.PasswordHasher
}

// @inject
func NewDeveloperUserBasicAuthMiddleware(userRepo repository.UserRepository, hasher security.PasswordHasher) *DeveloperUserBasicAuthMiddleware {
	return &DeveloperUserBasicAuthMiddleware{userRepo: userRepo, hasher: hasher}
}

func (this *DeveloperUserBasicAuthMiddleware) Wrap(next echo.HandlerFunc) echo.HandlerFunc {
	return middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		user, passwordHash, err := this.userRepo.GetByUsernameWithPassword(c.Request().Context(), username)
		if err != nil {
			return false, nil
		}
		if passwordHash == nil || !this.hasher.Equal(*passwordHash, password) {
			return false, nil
		}
		return user.Role == enum.RoleDeveloper, nil
	})(next)
}
