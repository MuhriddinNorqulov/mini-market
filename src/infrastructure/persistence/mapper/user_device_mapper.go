package mapper

import (
	"mini-market/src/core/domain/entity"
	"mini-market/src/infrastructure/persistence/models"
)

func UserDeviceModelToEntity(it *models.UserDeviceModel) *entity.UserDeviceEntity {
	if it == nil {
		return nil
	}
	return &entity.UserDeviceEntity{
		ID:             it.ID,
		UserID:         it.UserID,
		DeviceID:       it.DeviceID,
		Name:           it.Name,
		BrowserFamily:  it.BrowserFamily,
		PlatformFamily: it.PlatformFamily,
		LastSeenAt:     it.LastSeenAt,
		LastIP:         it.LastIP,
		TrustedAt:      it.TrustedAt,
		TrustExpiresAt: it.TrustExpiresAt,
	}
}

func UserDeviceEntityToModel(it *entity.UserDeviceEntity) *models.UserDeviceModel {
	return &models.UserDeviceModel{
		UserID:         it.UserID,
		DeviceID:       it.DeviceID,
		Name:           it.Name,
		BrowserFamily:  it.BrowserFamily,
		PlatformFamily: it.PlatformFamily,
		LastSeenAt:     it.LastSeenAt,
		LastIP:         it.LastIP,
		TrustedAt:      it.TrustedAt,
		TrustExpiresAt: it.TrustExpiresAt,
	}
}
