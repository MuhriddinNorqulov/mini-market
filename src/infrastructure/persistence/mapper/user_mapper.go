package mapper

import (
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
)

type UserRow struct {
	ID        uint
	Username  *string
	FirstName string
	LastName  *string
	Role      enum.Role
}

func UserRowToEntity(r *UserRow) *entity.UserEntity {
	return &entity.UserEntity{
		ID:        r.ID,
		Username:  r.Username,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Role:      r.Role,
	}
}
