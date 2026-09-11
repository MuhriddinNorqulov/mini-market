package groups

import (
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/entrypoint/http/handlers/order"
	"mini-market/src/entrypoint/http/interceptor/middlewares"
	"mini-market/src/entrypoint/http/interceptor/permissions"
)

type OrderGroup struct {
	createOrderHandler  *order.CreateOrderHandler
	getOrderHandler     *order.GetOrderHandler
	listOrdersHandler   *order.ListOrdersHandler
	cancelOrderHandler  *order.CancelOrderHandler
	confirmOrderHandler *order.ConfirmOrderHandler
}

// @inject
func NewOrderGroup(
	createOrderHandler *order.CreateOrderHandler,
	getOrderHandler *order.GetOrderHandler,
	listOrdersHandler *order.ListOrdersHandler,
	cancelOrderHandler *order.CancelOrderHandler,
	confirmOrderHandler *order.ConfirmOrderHandler,
) *OrderGroup {
	return &OrderGroup{
		createOrderHandler:  createOrderHandler,
		getOrderHandler:     getOrderHandler,
		listOrdersHandler:   listOrdersHandler,
		cancelOrderHandler:  cancelOrderHandler,
		confirmOrderHandler: confirmOrderHandler,
	}
}

func (this *OrderGroup) RegisterRoutes(g httpport.Group) {
	g.POST("", this.createOrderHandler.Handle, permissions.AuthenticatedUserPermission, middlewares.IdempotencyKeyMiddleware)
	g.GET("", this.listOrdersHandler.Handle, permissions.RolePermission(enum.RoleAdmin))
	g.GET("/:id", this.getOrderHandler.Handle, permissions.AuthenticatedUserPermission)
	g.POST("/:id/cancel", this.cancelOrderHandler.Handle, permissions.AuthenticatedUserPermission)
	g.POST("/:id/confirm", this.confirmOrderHandler.Handle, permissions.RolePermission(enum.RoleAdmin))
}
