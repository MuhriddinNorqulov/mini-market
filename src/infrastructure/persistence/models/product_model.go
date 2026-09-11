package models

import "gorm.io/gorm"

type ProductModel struct {
	gorm.Model

	Name  string `gorm:"size: 255;not null;"`
	Price int64  `gorm:"not null;"`

	StockQuantity int64 `gorm:"not null;min:0;"`
}

func (this *ProductModel) TableName() string {
	return "products"
}
