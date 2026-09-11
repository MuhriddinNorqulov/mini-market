package context

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/ports/httpport/request"
	"mini-market/src/infrastructure/echohttp/requestimpl"
	"net/http"

	"strings"

	"github.com/labstack/echo/v4"
)

const wrapperContextKey = "anticopy.ctx.EchoContext"

type EchoContext struct {
	user           *entity.UserEntity
	IdempotencyKey string

	echo.Context
}

func NewEchoContext(c echo.Context) ctx.Context {
	return &EchoContext{Context: c}
}

func (e *EchoContext) super() echo.Context {
	return e.Context
}

func (e *EchoContext) Store() {
	e.Set(wrapperContextKey, e)
}

func Lookup(c echo.Context) (*EchoContext, bool) {
	v := c.Get(wrapperContextKey)
	wrapped, ok := v.(*EchoContext)
	return wrapped, ok
}

func (e *EchoContext) Unwrap() echo.Context {
	return e.Context
}

func (e *EchoContext) GetContext() context.Context {
	return e.Context.Request().Context()
}

func (e *EchoContext) SetContext(ctx context.Context) {
	e.SetRequest(e.Context.Request().WithContext(ctx))
}

func (e *EchoContext) Param(name string) string {
	return e.Context.Param(name)
}

func (e *EchoContext) QueryParam(name string) string {
	return e.Context.QueryParam(name)
}

func (e *EchoContext) SetUser(user *entity.UserEntity) {
	e.user = user
}

func (e *EchoContext) User() *entity.UserEntity {
	return e.user
}

func (e *EchoContext) JSON(status int, i any) error {
	return e.super().JSON(status, i)
}

func (e *EchoContext) Success(status int, i any) error {
	payload := response.NewResponse(response.CodeSuccess, true, i, "")
	return e.JSON(status, payload)
}

func (e *EchoContext) GetRequest() *request.Request {
	header := e.super().Request().Header
	h := requestimpl.NewHeaderImpl(header)
	return &request.Request{Header: h}
}

func (e *EchoContext) GetIdempotencyKey() string {
	return e.IdempotencyKey
}

func (e *EchoContext) SetIdempotencyKey(key string) {
	e.IdempotencyKey = key
}

func (e *EchoContext) StreamFile(filename, contentType string, r io.Reader) error {
	e.super().Response().Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	return e.super().Stream(200, contentType, r)
}

func (e *EchoContext) GetBasicCredential() (string, string, error) {
	auth := e.super().Request().Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Basic ") {
		return "", "", errors.New("authorization header missing or invalid")
	}

	payload := strings.TrimPrefix(auth, "Basic ")

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", "", errors.New("invalid authorization header")
	}

	pair := strings.SplitN(string(decoded), ":", 2)
	if len(pair) != 2 {
		return "", "", errors.New("invalid authorization header")
	}

	return pair[0], pair[1], nil
}

func (e *EchoContext) Cookie(name string) (string, bool) {
	c, err := e.super().Cookie(name)
	if err != nil || c == nil || c.Value == "" {
		return "", false
	}
	return c.Value, true
}

func (e *EchoContext) SetCookie(cookie *ctx.Cookie) {
	e.super().SetCookie(&http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
		SameSite: toHTTPSameSite(cookie.SameSite),
	})
}

func toHTTPSameSite(s ctx.SameSite) http.SameSite {
	switch s {
	case ctx.SameSiteLaxMode:
		return http.SameSiteLaxMode
	case ctx.SameSiteStrictMode:
		return http.SameSiteStrictMode
	case ctx.SameSiteNoneMode:
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}
