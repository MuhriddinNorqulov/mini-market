package models

import (
	"mini-market/src/core/domain/entity/enum"
	"time"

	"gorm.io/gorm"
)

type OrderModel struct {
	gorm.Model

	// unique together idempotency_key + user_id
	IdempotencyKey string `gorm:"uniqueIndex:idx_order_idempotency_key;not null;size:64;"`

	UserID uint       `gorm:"not null;uniqueIndex:idx_order_idempotency_key;"`
	User   *UserModel `gorm:"constraint:OnDelete:SET NULL;"`

	TotalPrice int64 `gorm:"not null;"`

	Items []*OrderItemModel `gorm:"foreignKey:OrderID"`

	ExpiresAt    time.Time               `gorm:"not null;"`
	Status       enum.OrderStatus        `gorm:"not null;size:32;default:'pending';check:chk_orders_status,status IN ('pending','confirmed','cancelled')"`
	CancelReason *enum.OrderCancelReason `gorm:"size:32;check:chk_orders_cancel_reason,cancel_reason IN ('user_requested','timeout')"`
	CompletedAt  *time.Time
	CancelledAt  *time.Time
}

func (this *OrderModel) TableName() string {
	return "orders"
}
