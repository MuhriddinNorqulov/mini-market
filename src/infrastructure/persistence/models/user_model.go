package models

import (
	"time"

	"mini-market/src/core/domain/entity/enum"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	Username  *string `gorm:"uniqueIndex:idx_users_username;size:32;"`
	FirstName string  `gorm:"size:64;"`
	LastName  *string `gorm:"size:64;"`

	Role              enum.Role `gorm:"not null;default:'USER'"`
	Password          *string
	PasswordUpdatedAt *time.Time
}

func (this *UserModel) TableName() string {
	return "users"
}
