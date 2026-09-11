package middlewares

import (
	"errors"
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/ports/httpport"
	"mini-market/src/core/domain/ports/httpport/ctx"

	"github.com/labstack/echo/v4"
)

type ResponseMiddleware struct{}

// @inject
func NewResponseMiddleware() *ResponseMiddleware {
	return &ResponseMiddleware{}
}

func (this *ResponseMiddleware) Wrap(next httpport.HandlerFunc) httpport.HandlerFunc {
	return func(c ctx.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		if he, ok := errors.AsType[*echo.HTTPError](err); ok {
			return he
		}

		if resp, ok := errors.AsType[*response.Response](err); ok {
			return c.JSON(this.codeToStatus(resp.Code), resp)
		}

		return err
	}
}

func (this *ResponseMiddleware) codeToStatus(c response.Code) int {
	switch c {
	case response.CodeBadRequest:
		return 400
	case response.CodeUnauthorized:
		return 401
	case response.CodeExpiredToken:
		return 401
	case response.CodeInvalidToken:
		return 401
	case response.CodeForbidden:
		return 403
	case response.CodNotFound:
		return 404
	case response.CodeConflict:
		return 409
	case response.CodeInsufficientStock:
		return 409
	case response.CodeInvalidOrderStatus:
		return 422
	case response.CodeDatabaseError:
		return 500
	case response.CodeGatewayError:
		return 500
	}
	return 200
}
