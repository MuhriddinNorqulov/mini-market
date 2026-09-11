package entity

type ProductEntity struct {
	ID uint `json:"id"`

	Name  string `json:"name"`
	Price int64  `json:"price"`

	StockQuantity int64 `json:"stock_quantity"`
}
