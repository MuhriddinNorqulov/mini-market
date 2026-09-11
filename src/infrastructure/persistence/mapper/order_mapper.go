package mapper

import (
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"time"
)

type OrderRow struct {
	ID             uint
	UserID         uint
	IdempotencyKey string
	TotalPrice     int64
	Status         enum.OrderStatus
	ExpiresAt      time.Time
	CreatedAt      time.Time
}

type OrderItemRow struct {
	ID        uint
	ProductID uint
	Quantity  int64
	UnitPrice int64
}

func OrderRowToEntity(r *OrderRow) *entity.OrderEntity {
	return &entity.OrderEntity{
		ID:             r.ID,
		UserID:         r.UserID,
		IdempotencyKey: r.IdempotencyKey,
		TotalPrice:     r.TotalPrice,
		Status:         r.Status,
		ExpiresAt:      r.ExpiresAt,
		CreatedAt:      r.CreatedAt,
	}
}

func OrderItemRowToEntity(r *OrderItemRow) *entity.OrderItemEntity {
	return &entity.OrderItemEntity{
		ID:        r.ID,
		ProductID: r.ProductID,
		Quantity:  r.Quantity,
		UnitPrice: r.UnitPrice,
	}
}
