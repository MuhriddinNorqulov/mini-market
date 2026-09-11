package permissions_test

import (
	"context"
	"io"
	"mime/multipart"
	"testing"

	"mini-market/src/core/application/response"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/httpport/ctx"
	"mini-market/src/core/domain/ports/httpport/request"
	"mini-market/src/entrypoint/http/interceptor/permissions"
)

type fakeContext struct {
	ctx  context.Context
	user *entity.UserEntity
}

func (this *fakeContext) GetContext() context.Context     { return this.ctx }
func (this *fakeContext) SetContext(ctx context.Context)  { this.ctx = ctx }
func (this *fakeContext) User() *entity.UserEntity        { return this.user }
func (this *fakeContext) SetUser(user *entity.UserEntity) { this.user = user }
func (this *fakeContext) GetIdempotencyKey() string       { return "" }
func (this *fakeContext) SetIdempotencyKey(key string)    {}
func (this *fakeContext) JSON(status int, i any) error    { return nil }
func (this *fakeContext) Success(status int, i any) error { return nil }
func (this *fakeContext) StreamFile(filename, contentType string, r io.Reader) error {
	return nil
}
func (this *fakeContext) Param(name string) string      { return "" }
func (this *fakeContext) QueryParam(name string) string { return "" }
func (this *fakeContext) FormFile(name string) (*multipart.FileHeader, error) {
	return nil, nil
}
func (this *fakeContext) Bind(i any) error     { return nil }
func (this *fakeContext) Validate(i any) error { return nil }
func (this *fakeContext) GetRequest() *request.Request {
	return &request.Request{}
}
func (this *fakeContext) Cookie(name string) (string, bool) { return "", false }
func (this *fakeContext) SetCookie(cookie *ctx.Cookie)      {}
func (this *fakeContext) GetBasicCredential() (string, string, error) {
	return "", "", nil
}

func newFakeContext(user *entity.UserEntity) *fakeContext {
	return &fakeContext{ctx: context.Background(), user: user}
}

func TestRolePermission_AllowedRolePasses(t *testing.T) {
	next := func(c ctx.Context) error { return nil }
	mw := permissions.RolePermission(enum.RoleAdmin)

	err := mw(next)(newFakeContext(&entity.UserEntity{ID: 1, Role: enum.RoleAdmin}))
	if err != nil {
		t.Fatalf("expected no error for allowed role, got %v", err)
	}
}

func TestRolePermission_DisallowedRoleIsForbidden(t *testing.T) {
	next := func(c ctx.Context) error { return nil }
	mw := permissions.RolePermission(enum.RoleAdmin)

	err := mw(next)(newFakeContext(&entity.UserEntity{ID: 1, Role: enum.RoleUser}))
	if !response.IsErrorCode(err, response.CodeForbidden) {
		t.Fatalf("expected CodeForbidden, got %v", err)
	}
}

func TestRolePermission_NilUserIsUnauthorized(t *testing.T) {
	next := func(c ctx.Context) error { return nil }
	mw := permissions.RolePermission(enum.RoleAdmin)

	err := mw(next)(newFakeContext(nil))
	if !response.IsErrorCode(err, response.CodeUnauthorized) {
		t.Fatalf("expected CodeUnauthorized, got %v", err)
	}
}
