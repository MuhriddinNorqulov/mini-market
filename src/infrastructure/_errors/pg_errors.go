package _errors

import (
	"errors"

	"mini-market/src/core/application/response"
	"mini-market/src/core/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const pgUniqueViolation = "23505"

func RawSQLErrorWrap(err error) error {
	if err == nil {
		return nil
	}

	caller := utils.CallerPath(2)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return response.NewSafeError(response.CodeConflict, err, caller)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.NewSafeError(response.CodNotFound, err, caller)
	}

	return response.NewSafeError(response.CodeDatabaseError, err, caller)
}
