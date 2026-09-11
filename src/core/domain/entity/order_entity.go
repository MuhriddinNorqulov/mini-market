package entity

import (
	"mini-market/src/core/domain/entity/enum"
	"time"
)

type OrderEntity struct {
	ID uint `json:"id"`

	UserID         uint   `json:"user_id"`
	IdempotencyKey string `json:"idempotency_key"`

	TotalPrice int64 `json:"total_price"`

	Items []*OrderItemEntity `json:"items"`

	Status       enum.OrderStatus        `json:"status"`
	CancelReason *enum.OrderCancelReason `json:"cancel_reason,omitempty"`

	ExpiresAt   time.Time  `json:"expires_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
