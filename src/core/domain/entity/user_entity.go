package entity

import "mini-market/src/core/domain/entity/enum"

type UserEntity struct {
	ID uint `json:"id"`

	GoogleID      *string `json:"google_id"`
	Email         *string `json:"email"`
	Picture       *string `json:"picture"`
	EmailVerified bool    `json:"email_verified"`

	PhoneNumber *string `json:"phone_number"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	MiddleName  *string `json:"middle_name,omitempty"`

	ProfileImageFileID *uint `json:"-"`

	Role enum.Role `json:"role"`
}
