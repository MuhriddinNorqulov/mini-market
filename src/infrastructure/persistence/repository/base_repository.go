package repository

import (
	"context"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/infrastructure/db"

	"gorm.io/gorm"
)

type BaseRepository struct {
	database *db.Database
}

// @inject
func NewBaseRepository(database *db.Database) *BaseRepository {
	return &BaseRepository{database: database}
}

func (this *BaseRepository) db(ctx context.Context) *gorm.DB {
	return this.database.WithContext(ctx)
}

func (this *BaseRepository) tx(ctx context.Context, tx unitofwork.Tx) *gorm.DB {
	return tx.(*gorm.DB).WithContext(ctx)
}
