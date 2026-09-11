package order

import (
	"net/http"

	"mini-market/src/core/application/usecases/orderusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/query"
)

type ListMyOrdersHandler struct {
	uc *orderusecases.ListMyOrdersUseCase
}

// @inject
func NewListMyOrdersHandler(uc *orderusecases.ListMyOrdersUseCase) *ListMyOrdersHandler {
	return &ListMyOrdersHandler{uc: uc}
}

// Handle godoc
// @Tags         Orders
// @Summary      List my orders
// @Description  Returns a paginated list of the authenticated caller's own orders, newest first. The default page (page=1, default limit) is served from Redis (per-user cache key); any other page/limit combination reads straight from Postgres.
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     integer false  "Page number (default 1)"
// @Param        limit  query     integer false  "Items per page (default 20)"
// @Success      200    {object}  response.Response  "payload is a paging envelope: {page, size, total, items[]entity.OrderEntity}"
// @Router       /me/orders [get]
func (this *ListMyOrdersHandler) Handle(c ctx.Context) error {
	filter := query.OrderListFilter{
		Page:  ctx.GetIntQueryParam(c, "page", 1),
		Limit: ctx.GetIntQueryParam(c, "limit", query.DefaultPageLimit),
	}

	result, err := this.uc.Invoke(c.GetContext(), c.User(), filter)
	if err != nil {
		return err
	}
	return c.Success(http.StatusOK, result)
}
