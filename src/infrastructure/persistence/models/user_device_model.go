package models

import (
	"time"

	"gorm.io/gorm"
)

type UserDeviceModel struct {
	gorm.Model

	UserID uint       `gorm:"not null;uniqueIndex:idx_user_devices_user_device;"`
	User   *UserModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	DeviceID string `gorm:"size:64;not null;uniqueIndex:idx_user_devices_user_device;"`

	Name string `gorm:"size:128;not null;"`

	BrowserFamily  string `gorm:"size:32;not null;"`
	PlatformFamily string `gorm:"size:32;not null;"`

	LastSeenAt time.Time `gorm:"not null;"`
	LastIP     string    `gorm:"size:64;"`

	TrustedAt      *time.Time
	TrustExpiresAt *time.Time
}

func (this *UserDeviceModel) TableName() string {
	return "user_devices"
}
