package dto

type CreateOrderItemInput struct {
	ProductID uint  `json:"product_id" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderInput struct {
	Items []CreateOrderItemInput `json:"items" validate:"required,min=1,dive"`
}
