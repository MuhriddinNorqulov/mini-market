package repository

import (
	"context"
	"time"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/infrastructure/_errors"
	"mini-market/src/infrastructure/persistence/mapper"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	*BaseRepository
}

// @inject
func NewUserRepositoryImpl(baseRepository *BaseRepository) repository.UserRepository {
	return &UserRepositoryImpl{BaseRepository: baseRepository}
}

func (this *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.UserEntity, error) {
	db := this.db(ctx)
	var row mapper.UserRow
	err := db.Raw(
		`SELECT id, username, first_name, last_name, role
		 FROM users WHERE id = ? AND deleted_at IS NULL`,
		id,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return mapper.UserRowToEntity(&row), nil
}

func (this *UserRepositoryImpl) CreateTx(ctx context.Context, tx unitofwork.Tx, user *entity.UserEntity) (*entity.UserEntity, error) {
	db := this.tx(ctx, tx)
	var row mapper.UserRow
	err := db.Raw(
		`INSERT INTO users (username, first_name, last_name, role, created_at, updated_at)
		 VALUES (?, ?, ?, ?, now(), now())
		 RETURNING id, username, first_name, last_name, role`,
		user.Username, user.FirstName, user.LastName, user.Role,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	return mapper.UserRowToEntity(&row), nil
}

func (this *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.UserEntity, error) {
	db := this.db(ctx)
	var row mapper.UserRow
	err := db.Raw(
		`SELECT id, username, first_name, last_name, role
		 FROM users WHERE username = ? AND deleted_at IS NULL`,
		username,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return mapper.UserRowToEntity(&row), nil
}

func (this *UserRepositoryImpl) GetByUsernameWithPassword(ctx context.Context, username string) (*entity.UserEntity, *string, error) {
	db := this.db(ctx)
	var row struct {
		mapper.UserRow
		Password *string
	}
	err := db.Raw(
		`SELECT id, username, first_name, last_name, role, password
		 FROM users WHERE username = ? AND deleted_at IS NULL`,
		username,
	).Scan(&row).Error
	if err != nil {
		return nil, nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return mapper.UserRowToEntity(&row.UserRow), row.Password, nil
}

func (this *UserRepositoryImpl) GetByIDWithLock(ctx context.Context, tx unitofwork.Tx, id uint) (*entity.UserEntity, error) {
	db := this.tx(ctx, tx)
	var row mapper.UserRow
	err := db.Raw(
		`SELECT id, username, first_name, last_name, role
		 FROM users WHERE id = ? AND deleted_at IS NULL FOR UPDATE`,
		id,
	).Scan(&row).Error
	if err != nil {
		return nil, _errors.RawSQLErrorWrap(err)
	}
	if row.ID == 0 {
		return nil, _errors.RawSQLErrorWrap(gorm.ErrRecordNotFound)
	}
	return mapper.UserRowToEntity(&row), nil
}

func (this *UserRepositoryImpl) UpdatePasswordTx(ctx context.Context, tx unitofwork.Tx, userID uint, passwordHash string, updatedAt time.Time) error {
	result := this.tx(ctx, tx).Exec(
		`UPDATE users SET password = ?, password_updated_at = ?, updated_at = now() WHERE id = ?`,
		passwordHash, updatedAt, userID,
	)
	return _errors.RawSQLErrorWrap(result.Error)
}
