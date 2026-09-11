package product

import (
	"net/http"

	"mini-market/src/core/application/usecases/productusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/query"
)

type ListProductsHandler struct {
	uc *productusecases.ListProductsUseCase
}

// @inject
func NewListProductsHandler(uc *productusecases.ListProductsUseCase) *ListProductsHandler {
	return &ListProductsHandler{uc: uc}
}

// Handle godoc
// @Tags         Products
// @Summary      List products
// @Description  Returns a paginated product list, ordered by id. The default page (page=1, default limit) is served from Redis; any other page/limit combination reads straight from Postgres.
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     integer false  "Page number (default 1)"
// @Param        limit  query     integer false  "Items per page (default 20)"
// @Success      200    {object}  response.Response  "payload is a paging envelope: {page, size, total, items[]entity.ProductEntity}"
// @Router       /products [get]
func (this *ListProductsHandler) Handle(c ctx.Context) error {
	filter := query.ProductListFilter{
		Page:  ctx.GetIntQueryParam(c, "page", 1),
		Limit: ctx.GetIntQueryParam(c, "limit", query.DefaultPageLimit),
	}

	result, err := this.uc.Invoke(c.GetContext(), filter)
	if err != nil {
		return err
	}
	return c.Success(http.StatusOK, result)
}
