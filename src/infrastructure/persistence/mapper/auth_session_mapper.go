package mapper

import (
	"mini-market/src/core/domain/entity"
	"mini-market/src/infrastructure/persistence/models"
)

func AuthSessionModelToEntity(it *models.AuthSessionModel) *entity.AuthSessionEntity {
	if it == nil {
		return nil
	}
	return &entity.AuthSessionEntity{
		ID:               it.ID,
		SessionID:        it.SessionID,
		UserID:           it.UserID,
		UserDeviceID:     it.UserDeviceID,
		RefreshTokenHash: it.RefreshTokenHash,
		IssuedAt:         it.IssuedAt,
		LastSeenAt:       it.LastSeenAt,
		ExpiresAt:        it.ExpiresAt,
		RevokedAt:        it.RevokedAt,
		RevokedReason:    it.RevokedReason,
		IP:               it.IP,
		UserAgent:        it.UserAgent,
		Device:           UserDeviceModelToEntity(it.UserDevice),
	}
}

func AuthSessionEntityToModel(it *entity.AuthSessionEntity) *models.AuthSessionModel {
	return &models.AuthSessionModel{
		SessionID:        it.SessionID,
		UserID:           it.UserID,
		UserDeviceID:     it.UserDeviceID,
		RefreshTokenHash: it.RefreshTokenHash,
		IssuedAt:         it.IssuedAt,
		LastSeenAt:       it.LastSeenAt,
		ExpiresAt:        it.ExpiresAt,
		RevokedAt:        it.RevokedAt,
		RevokedReason:    it.RevokedReason,
		IP:               it.IP,
		UserAgent:        it.UserAgent,
	}
}
