package product

import (
	"net/http"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/usecases/productusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type CreateProductHandler struct {
	uc *productusecases.CreateProductUseCase
}

// @inject
func NewCreateProductHandler(uc *productusecases.CreateProductUseCase) *CreateProductHandler {
	return &CreateProductHandler{uc: uc}
}

// Handle godoc
// @Tags         Products
// @Summary      Create a product (admin)
// @Description  Creates a product with the given name, price, and initial stock quantity. Admin-only.
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateProductInput  true  "Product name, price, and stock quantity"
// @Success      201   {object}  response.Response{payload=entity.ProductEntity}
// @Failure      400   {object}  response.Response  "BAD_REQUEST — validation error"
// @Failure      401   {object}  response.Response  "UNAUTHORIZED — missing or invalid token"
// @Failure      403   {object}  response.Response  "FORBIDDEN — admin role required"
// @Router       /products [post]
func (this *CreateProductHandler) Handle(c ctx.Context) error {
	input, err := ctx.GetBody[dto.CreateProductInput](c)
	if err != nil {
		return err
	}

	product, err := this.uc.Invoke(c.GetContext(), *input)
	if err != nil {
		return err
	}
	return c.Success(http.StatusCreated, product)
}
