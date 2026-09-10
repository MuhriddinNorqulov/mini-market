package repository

import (
	"context"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/unitofwork"
	"time"
)

type AuthSessionRepository interface {
	CreateTx(ctx context.Context, tx unitofwork.Tx, session *entity.AuthSessionEntity) (*entity.AuthSessionEntity, error)

	GetBySessionID(ctx context.Context, sessionID string) (*entity.AuthSessionEntity, error)

	ListActiveByUserID(ctx context.Context, userID uint, now time.Time) ([]*entity.AuthSessionEntity, error)

	RotateTx(ctx context.Context, tx unitofwork.Tx, sessionID, newHash string, lastSeenAt, expiresAt time.Time) error

	RevokeBySessionIDTx(ctx context.Context, tx unitofwork.Tx, userID uint, sessionID string, reason enum.SessionRevokeReason, at time.Time) (bool, error)

	RevokeActiveByUserDeviceTx(ctx context.Context, tx unitofwork.Tx, userID, userDeviceID uint, reason enum.SessionRevokeReason, at time.Time) ([]string, error)

	RevokeAllByUserIDTx(ctx context.Context, tx unitofwork.Tx, userID uint, exceptSessionID string, reason enum.SessionRevokeReason, at time.Time) ([]string, error)
}
