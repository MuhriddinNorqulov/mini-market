package entity

import "time"

type UserDeviceEntity struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`

	DeviceID string `json:"device_id"`

	Name           string `json:"name"`
	BrowserFamily  string `json:"browser_family"`
	PlatformFamily string `json:"platform_family"`

	LastSeenAt time.Time `json:"last_seen_at"`
	LastIP     string    `json:"last_ip"`

	TrustedAt      *time.Time `json:"trusted_at"`
	TrustExpiresAt *time.Time `json:"trust_expires_at"`
}
