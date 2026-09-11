package dto

type CreateProductInput struct {
	Name          string `json:"name" validate:"required,min=1,max=255"`
	Price         int64  `json:"price" validate:"required,gt=0"`
	StockQuantity int64  `json:"stock_quantity" validate:"gte=0"`
}
