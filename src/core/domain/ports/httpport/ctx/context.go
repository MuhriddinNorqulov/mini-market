package ctx

import (
	"context"
	"io"
	"mime/multipart"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/ports/httpport/request"
)

type Context interface {
	GetContext() context.Context

	SetContext(ctx context.Context)
	User() *entity.UserEntity
	SetUser(user *entity.UserEntity)

	GetIdempotencyKey() string
	SetIdempotencyKey(key string)

	JSON(status int, i any) error
	Success(status int, i any) error
	StreamFile(filename, contentType string, r io.Reader) error

	Param(name string) string

	QueryParam(name string) string

	FormFile(name string) (*multipart.FileHeader, error)

	Bind(i any) error
	Validate(i any) error

	GetRequest() *request.Request

	Cookie(name string) (string, bool)
	SetCookie(cookie *Cookie)

	GetBasicCredential() (string, string, error)
}
