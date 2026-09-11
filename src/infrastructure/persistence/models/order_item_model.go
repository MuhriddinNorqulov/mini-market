package models

import "gorm.io/gorm"

type OrderItemModel struct {
	gorm.Model

	// unique together order_id + product_id
	OrderID uint        `gorm:"not null;uniqueIndex:idx_order_items_order_id_product_id;"`
	Order   *OrderModel `gorm:"constraint:OnDelete:SET NULL;"`

	ProductID uint          `gorm:"not null;uniqueIndex:idx_order_items_order_id_product_id;"`
	Product   *ProductModel `gorm:"constraint:OnDelete:SET NULL;"`

	Quantity  int64 `gorm:"not null;"`
	UnitPrice int64 `gorm:"not null;"`
}

func (this *OrderItemModel) TableName() string {
	return "order_items"
}
