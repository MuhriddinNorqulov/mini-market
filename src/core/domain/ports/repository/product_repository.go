package repository

import (
	"context"

	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/unitofwork"
	"mini-market/src/core/domain/query"
)

type StockReservation struct {
	ProductID uint
	Quantity  int64
}

type ProductRepository interface {
	Create(ctx context.Context, product *entity.ProductEntity) (*entity.ProductEntity, error)
	GetByID(ctx context.Context, id uint) (*entity.ProductEntity, error)
	GetByIDs(ctx context.Context, ids []uint) ([]*entity.ProductEntity, error)
	List(ctx context.Context, filter query.ProductListFilter) (*entity.PagingEntity[*entity.ProductEntity], error)
	ReserveStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []StockReservation) (prices map[uint]int64, failedProductIDs []uint, err error)
	ReleaseStockBatchTx(ctx context.Context, tx unitofwork.Tx, reservations []StockReservation) error
}
