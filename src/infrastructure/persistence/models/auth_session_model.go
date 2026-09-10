package models

import (
	"mini-market/src/core/domain/entity/enum"
	"time"

	"gorm.io/gorm"
)

type AuthSessionModel struct {
	gorm.Model

	SessionID string `gorm:"size:64;uniqueIndex;not null;"`

	UserID uint       `gorm:"not null;index:idx_auth_sessions_user_active,priority:1;"`
	User   *UserModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserDeviceID uint             `gorm:"not null;index;"`
	UserDevice   *UserDeviceModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	RefreshTokenHash string `gorm:"size:64;uniqueIndex;not null;"`

	IssuedAt   time.Time `gorm:"not null;"`
	LastSeenAt time.Time `gorm:"not null;"`
	ExpiresAt  time.Time `gorm:"not null;index:idx_auth_sessions_user_active,priority:3;"`

	RevokedAt     *time.Time                `gorm:"index:idx_auth_sessions_user_active,priority:2;"`
	RevokedReason *enum.SessionRevokeReason `gorm:"size:32;"`

	IP        string `gorm:"size:64;"`
	UserAgent string `gorm:"size:512;"`
}

func (this *AuthSessionModel) TableName() string {
	return "auth_sessions"
}
