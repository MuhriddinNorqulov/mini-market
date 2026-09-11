package mapper

import "mini-market/src/core/domain/entity"

type ProductRow struct {
	ID            uint
	Name          string
	Price         int64
	StockQuantity int64
}

type ReservedProductRow struct {
	ID    uint
	Price int64
}

func ProductRowToEntity(r *ProductRow) *entity.ProductEntity {
	return &entity.ProductEntity{
		ID:            r.ID,
		Name:          r.Name,
		Price:         r.Price,
		StockQuantity: r.StockQuantity,
	}
}
