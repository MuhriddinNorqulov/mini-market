package groups

import (
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/entrypoint/http/handlers/product"
	"mini-market/src/entrypoint/http/interceptor/permissions"
)

type ProductGroup struct {
	createProductHandler *product.CreateProductHandler
	getProductHandler    *product.GetProductHandler
	listProductsHandler  *product.ListProductsHandler
}

// @inject
func NewProductGroup(
	createProductHandler *product.CreateProductHandler,
	getProductHandler *product.GetProductHandler,
	listProductsHandler *product.ListProductsHandler,
) *ProductGroup {
	return &ProductGroup{
		createProductHandler: createProductHandler,
		getProductHandler:    getProductHandler,
		listProductsHandler:  listProductsHandler,
	}
}

func (this *ProductGroup) RegisterRoutes(g httpport.Group) {
	g.POST("", this.createProductHandler.Handle, permissions.RolePermission(enum.RoleAdmin))
	g.GET("", this.listProductsHandler.Handle, permissions.AuthenticatedUserPermission)
	g.GET("/:id", this.getProductHandler.Handle, permissions.AuthenticatedUserPermission)
}
