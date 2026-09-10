package _errors

import (
	"errors"
	"mini-market/src/core/application/response"
	"mini-market/src/core/utils"

	"gorm.io/gorm"
)

func GormErrorWrap(err error) error {
	if err == nil {
		return nil
	}

	caller := utils.CallerPath(2)

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return response.NewSafeError(response.CodNotFound, err, caller)

	}

	return response.NewSafeError(response.CodeDatabaseError, err, caller)
}
