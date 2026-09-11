package permissions

import (
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

func RolePermission(roles ...enum.Role) httpport.Middleware {
	allowed := make(map[enum.Role]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next httpport.HandlerFunc) httpport.HandlerFunc {
		return func(c ctx.Context) error {
			user := c.User()
			if user == nil {
				return response.NewResponse(response.CodeUnauthorized, false, nil, "Unauthorized user")
			}
			if !allowed[user.Role] {
				return response.NewResponse(response.CodeForbidden, false, nil, "insufficient role")
			}
			return next(c)
		}
	}
}
