package db

import (
	"mini-market/src/core/domain/ports/unitofwork"

	"gorm.io/gorm"
)

type AtomicTxImpl struct {
	db *Database
}

// @inject
func NewAtomicTxImpl(db *Database) unitofwork.Atomic {
	return &AtomicTxImpl{db: db}
}

func (this *AtomicTxImpl) Transaction(fn func(unitofwork.Tx) error) error {
	return this.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
