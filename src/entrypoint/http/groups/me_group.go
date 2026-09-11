package groups

import (
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/entrypoint/http/handlers/order"
	"mini-market/src/entrypoint/http/interceptor/permissions"
)

type MeGroup struct {
	listMyOrdersHandler *order.ListMyOrdersHandler
}

// @inject
func NewMeGroup(listMyOrdersHandler *order.ListMyOrdersHandler) *MeGroup {
	return &MeGroup{listMyOrdersHandler: listMyOrdersHandler}
}

func (this *MeGroup) RegisterRoutes(g httpport.Group) {
	g.GET("/orders", this.listMyOrdersHandler.Handle, permissions.AuthenticatedUserPermission)
}
