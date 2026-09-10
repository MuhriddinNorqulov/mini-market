package entity

import (
	"mini-market/src/core/domain/entity/enum"
	"time"
)

type AuthSessionEntity struct {
	ID uint `json:"id"`

	SessionID string `json:"session_id"`

	UserID       uint `json:"user_id"`
	UserDeviceID uint `json:"user_device_id"`

	RefreshTokenHash string `json:"-"`

	IssuedAt   time.Time `json:"issued_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	ExpiresAt  time.Time `json:"expires_at"`

	RevokedAt     *time.Time                `json:"revoked_at"`
	RevokedReason *enum.SessionRevokeReason `json:"revoked_reason"`

	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`

	Device *UserDeviceEntity `json:"device,omitempty"`
}

func (this *AuthSessionEntity) IsActive(now time.Time) bool {
	return this.RevokedAt == nil && this.ExpiresAt.After(now)
}
