package order

import (
	"net/http"

	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type ConfirmOrderHandler struct {
	uc *orderusecases.ConfirmOrderUseCase
}

// @inject
func NewConfirmOrderHandler(uc *orderusecases.ConfirmOrderUseCase) *ConfirmOrderHandler {
	return &ConfirmOrderHandler{uc: uc}
}

// Handle godoc
// @Tags         Orders
// @Summary      Confirm an order (admin)
// @Description  Confirms a pending order, standing in for a payment-captured webhook (no external payment provider is in scope). Admin-only.
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      integer  true  "Order ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response  "BAD_REQUEST — invalid id"
// @Failure      403  {object}  response.Response  "FORBIDDEN — admin role required"
// @Failure      404  {object}  response.Response  "NOT_FOUND — order does not exist"
// @Failure      422  {object}  response.Response  "INVALID_ORDER_STATUS — order is not pending"
// @Router       /orders/{id}/confirm [post]
func (this *ConfirmOrderHandler) Handle(c ctx.Context) error {
	id, err := ctx.GetIntPathParam(c, "id")
	if err != nil {
		return err
	}

	if err := this.uc.Invoke(c.GetContext(), uint(id)); err != nil {
		return err
	}
	return c.Success(http.StatusOK, nil)
}
