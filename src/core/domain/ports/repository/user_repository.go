package repository

import (
	"context"
	"time"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/unitofwork"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*entity.UserEntity, error)
	CreateTx(ctx context.Context, tx unitofwork.Tx, user *entity.UserEntity) (*entity.UserEntity, error)
	GetByUsernameWithPassword(ctx context.Context, username string) (*entity.UserEntity, *string, error)
	GetByID(ctx context.Context, id uint) (*entity.UserEntity, error)
	GetByIDWithLock(ctx context.Context, tx unitofwork.Tx, id uint) (*entity.UserEntity, error)
	UpdatePasswordTx(ctx context.Context, tx unitofwork.Tx, userID uint, passwordHash string, updatedAt time.Time) error
}
