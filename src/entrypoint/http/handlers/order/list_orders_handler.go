package order

import (
	"net/http"

	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/query"
)

type ListOrdersHandler struct {
	uc *orderusecases.ListOrdersUseCase
}

// @inject
func NewListOrdersHandler(uc *orderusecases.ListOrdersUseCase) *ListOrdersHandler {
	return &ListOrdersHandler{uc: uc}
}

// Handle godoc
// @Tags         Orders
// @Summary      List orders (admin)
// @Description  Returns a paginated list of the most recent orders across all users, newest first. Admin-only. The default page (page=1, default limit) is served from Redis; any other page/limit combination reads straight from Postgres.
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     integer false  "Page number (default 1)"
// @Param        limit  query     integer false  "Items per page (default 20)"
// @Success      200    {object}  response.Response  "payload is a paging envelope: {page, size, total, items[]entity.OrderEntity}"
// @Failure      403    {object}  response.Response  "FORBIDDEN — admin role required"
// @Router       /orders [get]
func (this *ListOrdersHandler) Handle(c ctx.Context) error {
	filter := query.OrderListFilter{
		Page:  ctx.GetIntQueryParam(c, "page", 1),
		Limit: ctx.GetIntQueryParam(c, "limit", query.DefaultPageLimit),
	}

	result, err := this.uc.Invoke(c.GetContext(), filter)
	if err != nil {
		return err
	}
	return c.Success(http.StatusOK, result)
}
