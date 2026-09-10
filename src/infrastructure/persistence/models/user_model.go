package models

import (
	"time"

	"mini-market/src/core/domain/entity/enum"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	GoogleID      *string `gorm:"uniqueIndex:idx_users_google_id;size:128;"`
	Email         *string `gorm:"uniqueIndex:idx_users_email;size:256;"`
	Picture       *string `gorm:"size:256;"`
	EmailVerified bool    `gorm:"not null;default:false;"`

	PhoneNumber *string `gorm:"uniqueIndex:idx_users_phone_number;size:32;"`
	FirstName   string  `gorm:"size:64;"`
	LastName    *string `gorm:"size:64;"`
	MiddleName  *string `gorm:"size:64;"`

	ProfileImageFileID *uint `gorm:"index;"`

	Role              enum.Role `gorm:"not null;default:'USER'"`
	Password          *string
	PasswordUpdatedAt *time.Time
}

func (this *UserModel) TableName() string {
	return "users"
}
