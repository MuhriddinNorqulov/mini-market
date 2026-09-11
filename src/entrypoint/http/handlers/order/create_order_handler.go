package order

import (
	"net/http"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type CreateOrderHandler struct {
	uc *orderusecases.CreateOrderUseCase
}

// @inject
func NewCreateOrderHandler(uc *orderusecases.CreateOrderUseCase) *CreateOrderHandler {
	return &CreateOrderHandler{uc: uc}
}

// Handle godoc
// @Tags         Orders
// @Summary      Create an order
// @Description  Atomically reserves stock for each item and creates a pending order that auto-cancels after 15 minutes if unpaid. Duplicate items in the same request are aggregated by quantity. Resubmitting the same Idempotency-Key returns the original order instead of reserving stock again.
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Idempotency-Key  header    string           true  "Client-generated key, unique per order attempt"
// @Param        body             body      dto.CreateOrderInput  true  "Order items (product_id, quantity)"
// @Success      201  {object}  response.Response{payload=entity.OrderEntity}
// @Failure      400  {object}  response.Response  "BAD_REQUEST — validation error or missing Idempotency-Key header"
// @Failure      404  {object}  response.Response  "NOT_FOUND — one or more product ids do not exist"
// @Failure      409  {object}  response.Response  "INSUFFICIENT_STOCK — not enough stock for one or more products"
// @Router       /orders [post]
func (this *CreateOrderHandler) Handle(c ctx.Context) error {
	input, err := ctx.GetBody[dto.CreateOrderInput](c)
	if err != nil {
		return err
	}

	order, err := this.uc.Invoke(c.GetContext(), c.User().ID, c.GetIdempotencyKey(), input.Items)
	if err != nil {
		return err
	}
	return c.Success(http.StatusCreated, order)
}
