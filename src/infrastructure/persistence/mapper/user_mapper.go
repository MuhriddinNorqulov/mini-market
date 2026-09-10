package mapper

import (
	"mini-market/src/core/domain/entity"
	"mini-market/src/infrastructure/persistence/models"
)

func UserModelToEntity(it *models.UserModel) *entity.UserEntity {
	return &entity.UserEntity{
		ID:                 it.ID,
		GoogleID:           it.GoogleID,
		Email:              it.Email,
		Picture:            it.Picture,
		EmailVerified:      it.EmailVerified,
		PhoneNumber:        it.PhoneNumber,
		FirstName:          it.FirstName,
		LastName:           it.LastName,
		MiddleName:         it.MiddleName,
		ProfileImageFileID: it.ProfileImageFileID,
		Role:               it.Role,
	}
}

func UserEntityToModel(it *entity.UserEntity) *models.UserModel {
	return &models.UserModel{
		GoogleID:           it.GoogleID,
		Email:              it.Email,
		Picture:            it.Picture,
		EmailVerified:      it.EmailVerified,
		PhoneNumber:        it.PhoneNumber,
		FirstName:          it.FirstName,
		LastName:           it.LastName,
		MiddleName:         it.MiddleName,
		ProfileImageFileID: it.ProfileImageFileID,
	}
}
