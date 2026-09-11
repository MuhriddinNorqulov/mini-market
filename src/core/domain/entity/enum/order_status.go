package enum

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderCancelReason string

const (
	OrderCancelReasonUserRequested OrderCancelReason = "user_requested"
	OrderCancelReasonTimeout       OrderCancelReason = "timeout"
)
