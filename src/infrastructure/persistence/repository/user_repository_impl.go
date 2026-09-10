package repository

import (
	"context"
	"time"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/infrastructure/_errors"
	"mini-market/src/infrastructure/persistence/mapper"
	"mini-market/src/infrastructure/persistence/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewUserRepositoryImpl(baseRepository *BaseRepository) repository.UserRepository {
	return &UserRepositoryImpl{BaseRepository: baseRepository}
}

func (this *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.UserEntity, error) {
	var _model models.UserModel
	if err := this.db(ctx).First(&_model, id).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.UserModelToEntity(&_model), nil
}

func (this *UserRepositoryImpl) CreateTx(ctx context.Context, tx unitofwork.Tx, user *entity.UserEntity) (*entity.UserEntity, error) {
	_model := mapper.UserEntityToModel(user)
	if err := this.tx(ctx, tx).Create(_model).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.UserModelToEntity(_model), nil
}

func (this *UserRepositoryImpl) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.UserEntity, error) {
	var _model models.UserModel
	if err := this.db(ctx).Model(&models.UserModel{}).Where("phone_number = ?", phoneNumber).First(&_model).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.UserModelToEntity(&_model), nil
}

func (this *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	var _model models.UserModel
	if err := this.db(ctx).Model(&models.UserModel{}).Where("email = ?", email).First(&_model).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.UserModelToEntity(&_model), nil
}

func (this *UserRepositoryImpl) GetByPhoneNumberWithPassword(ctx context.Context, phoneNumber string) (*entity.UserEntity, *string, error) {
	var _model models.UserModel
	if err := this.db(ctx).Model(&models.UserModel{}).Where("phone_number = ?", phoneNumber).First(&_model).Error; err != nil {
		return nil, nil, _errors.GormErrorWrap(err)
	}
	return mapper.UserModelToEntity(&_model), _model.Password, nil
}

func (this *UserRepositoryImpl) GetByIDWithLock(ctx context.Context, tx unitofwork.Tx, id uint) (*entity.UserEntity, error) {
	var model models.UserModel
	db := tx.(*gorm.DB).WithContext(ctx)
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&model, id).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.UserModelToEntity(&model), nil
}

func (this *UserRepositoryImpl) UpdatePasswordTx(ctx context.Context, tx unitofwork.Tx, userID uint, passwordHash string, updatedAt time.Time) error {
	return _errors.GormErrorWrap(
		this.tx(ctx, tx).Model(&models.UserModel{}).
			Where("id = ?", userID).
			Updates(map[string]any{"password": passwordHash, "password_updated_at": updatedAt}).Error,
	)
}
