package order

import (
	"net/http"

	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type CancelOrderHandler struct {
	uc *orderusecases.CancelOrderUseCase
}

// @inject
func NewCancelOrderHandler(uc *orderusecases.CancelOrderUseCase) *CancelOrderHandler {
	return &CancelOrderHandler{uc: uc}
}

// Handle godoc
// @Tags         Orders
// @Summary      Cancel an order
// @Description  Cancels a pending order and releases its reserved stock. Uses the same use case as the 15-minute auto-cancel background job.
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Order ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response  "BAD_REQUEST — invalid id"
// @Failure      403  {object}  response.Response  "FORBIDDEN — order belongs to another user"
// @Failure      404  {object}  response.Response  "NOT_FOUND — order does not exist"
// @Failure      422  {object}  response.Response  "INVALID_ORDER_STATUS — order is not pending"
// @Router       /orders/{id}/cancel [post]
func (this *CancelOrderHandler) Handle(c ctx.Context) error {
	id, err := ctx.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	if err := this.uc.Invoke(c.GetContext(), c.User(), uint(id), enum.OrderCancelReasonUserRequested); err != nil {
		return err
	}
	return c.Success(http.StatusOK, nil)
}
