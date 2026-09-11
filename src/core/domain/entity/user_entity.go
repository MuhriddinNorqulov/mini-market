package entity

import "mini-market/src/core/domain/entity/enum"

type UserEntity struct {
	ID uint `json:"id"`

	Username  *string `json:"username"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`

	Role enum.Role `json:"role"`
}
