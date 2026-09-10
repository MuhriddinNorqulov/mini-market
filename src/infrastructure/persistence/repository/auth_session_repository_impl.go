package repository

import (
	"context"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/infrastructure/_errors"
	"mini-market/src/infrastructure/persistence/mapper"
	"mini-market/src/infrastructure/persistence/models"
	"time"

	"gorm.io/gorm"
)

type AuthSessionRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewAuthSessionRepositoryImpl(baseRepository *BaseRepository) repository.AuthSessionRepository {
	return &AuthSessionRepositoryImpl{BaseRepository: baseRepository}
}

func (this *AuthSessionRepositoryImpl) CreateTx(ctx context.Context, tx unitofwork.Tx, session *entity.AuthSessionEntity) (*entity.AuthSessionEntity, error) {
	_model := mapper.AuthSessionEntityToModel(session)
	if err := this.tx(ctx, tx).Create(_model).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.AuthSessionModelToEntity(_model), nil
}

func (this *AuthSessionRepositoryImpl) GetBySessionID(ctx context.Context, sessionID string) (*entity.AuthSessionEntity, error) {
	var _model models.AuthSessionModel
	if err := this.db(ctx).Where("session_id = ?", sessionID).First(&_model).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	return mapper.AuthSessionModelToEntity(&_model), nil
}

func (this *AuthSessionRepositoryImpl) ListActiveByUserID(ctx context.Context, userID uint, now time.Time) ([]*entity.AuthSessionEntity, error) {
	var items []models.AuthSessionModel
	if err := this.db(ctx).
		Preload("UserDevice").
		Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, now).
		Order("last_seen_at DESC").
		Find(&items).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}

	result := make([]*entity.AuthSessionEntity, len(items))
	for i := range items {
		result[i] = mapper.AuthSessionModelToEntity(&items[i])
	}
	return result, nil
}

func (this *AuthSessionRepositoryImpl) RotateTx(ctx context.Context, tx unitofwork.Tx, sessionID, newHash string, lastSeenAt, expiresAt time.Time) error {
	return _errors.GormErrorWrap(
		this.tx(ctx, tx).Model(&models.AuthSessionModel{}).
			Where("session_id = ?", sessionID).
			Updates(map[string]any{
				"refresh_token_hash": newHash,
				"last_seen_at":       lastSeenAt,
				"expires_at":         expiresAt,
			}).Error,
	)
}

func (this *AuthSessionRepositoryImpl) RevokeBySessionIDTx(ctx context.Context, tx unitofwork.Tx, userID uint, sessionID string, reason enum.SessionRevokeReason, at time.Time) (bool, error) {

	res := this.tx(ctx, tx).Model(&models.AuthSessionModel{}).
		Where("session_id = ? AND user_id = ? AND revoked_at IS NULL", sessionID, userID).
		Updates(map[string]any{"revoked_at": at, "revoked_reason": reason})

	if res.Error != nil {
		return false, _errors.GormErrorWrap(res.Error)
	}
	return res.RowsAffected > 0, nil
}

func (this *AuthSessionRepositoryImpl) RevokeActiveByUserDeviceTx(ctx context.Context, tx unitofwork.Tx, userID, userDeviceID uint, reason enum.SessionRevokeReason, at time.Time) ([]string, error) {
	return this.revokeReturningSessionIDs(
		ctx, tx, reason, at,
		this.tx(ctx, tx).Model(&models.AuthSessionModel{}).
			Where("user_id = ? AND user_device_id = ? AND revoked_at IS NULL", userID, userDeviceID),
	)
}

func (this *AuthSessionRepositoryImpl) RevokeAllByUserIDTx(ctx context.Context, tx unitofwork.Tx, userID uint, exceptSessionID string, reason enum.SessionRevokeReason, at time.Time) ([]string, error) {
	q := this.tx(ctx, tx).Model(&models.AuthSessionModel{}).
		Where("user_id = ? AND revoked_at IS NULL", userID)

	if exceptSessionID != "" {
		q = q.Where("session_id <> ?", exceptSessionID)
	}

	return this.revokeReturningSessionIDs(ctx, tx, reason, at, q)
}

func (this *AuthSessionRepositoryImpl) revokeReturningSessionIDs(ctx context.Context, tx unitofwork.Tx, reason enum.SessionRevokeReason, at time.Time, q *gorm.DB) ([]string, error) {
	var sessionIDs []string
	if err := q.Session(&gorm.Session{}).Pluck("session_id", &sessionIDs).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}
	if len(sessionIDs) == 0 {
		return nil, nil
	}

	if err := this.tx(ctx, tx).Model(&models.AuthSessionModel{}).
		Where("session_id IN ?", sessionIDs).
		Updates(map[string]any{"revoked_at": at, "revoked_reason": reason}).Error; err != nil {
		return nil, _errors.GormErrorWrap(err)
	}

	return sessionIDs, nil
}
