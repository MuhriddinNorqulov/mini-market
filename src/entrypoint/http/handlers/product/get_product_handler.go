package product

import (
	"net/http"

	"mini-market/src/core/application/usecases/productusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type GetProductHandler struct {
	uc *productusecases.GetProductUseCase
}

// @inject
func NewGetProductHandler(uc *productusecases.GetProductUseCase) *GetProductHandler {
	return &GetProductHandler{uc: uc}
}

// Handle godoc
// @Tags         Products
// @Summary      Get a product by ID
// @Description  Returns a single product. Cached in Redis (cache-aside, invalidated on stock-changing writes).
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Product ID"
// @Success      200  {object}  response.Response{payload=entity.ProductEntity}
// @Failure      400  {object}  response.Response  "BAD_REQUEST — invalid id"
// @Failure      404  {object}  response.Response  "NOT_FOUND — product does not exist"
// @Router       /products/{id} [get]
func (this *GetProductHandler) Handle(c ctx.Context) error {
	id, err := ctx.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	product, err := this.uc.Invoke(c.GetContext(), uint(id))
	if err != nil {
		return err
	}
	return c.Success(http.StatusOK, product)
}
