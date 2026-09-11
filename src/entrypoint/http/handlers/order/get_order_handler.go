package order

import (
	"net/http"

	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type GetOrderHandler struct {
	uc *orderusecases.GetOrderUseCase
}

// @inject
func NewGetOrderHandler(uc *orderusecases.GetOrderUseCase) *GetOrderHandler {
	return &GetOrderHandler{uc: uc}
}

// Handle godoc
// @Tags         Orders
// @Summary      Get an order by ID
// @Description  Returns an order with its items and current status.
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Order ID"
// @Success      200  {object}  response.Response{payload=entity.OrderEntity}
// @Failure      400  {object}  response.Response  "BAD_REQUEST — invalid id"
// @Failure      403  {object}  response.Response  "FORBIDDEN — order belongs to another user"
// @Failure      404  {object}  response.Response  "NOT_FOUND — order does not exist"
// @Router       /orders/{id} [get]
func (this *GetOrderHandler) Handle(c ctx.Context) error {
	id, err := ctx.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	order, err := this.uc.Invoke(c.GetContext(), c.User(), uint(id))
	if err != nil {
		return err
	}
	return c.Success(http.StatusOK, order)
}
