package entity

type OrderItemEntity struct {
	ID uint `json:"id"`

	ProductID uint  `json:"product_id"`
	Quantity  int64 `json:"quantity"`
	UnitPrice int64 `json:"unit_price"`
}
