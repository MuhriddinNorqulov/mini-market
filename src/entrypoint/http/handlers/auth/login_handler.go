package auth

import (
	"net/http"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/usecases/authusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type LoginHandler struct {
	uc *authusecases.LoginUseCase
}

// @inject
func NewLoginHandler(uc *authusecases.LoginUseCase) *LoginHandler {
	return &LoginHandler{uc: uc}
}

// Handle godoc
// @Tags         Auth
// @Summary      Login
// @Description  Verifies a username and password and returns an access token.
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginInput  true  "Username and password"
// @Success      200   {object}  response.Response{payload=dto.AuthOutput}
// @Failure      400   {object}  response.Response  "BAD_REQUEST — validation error"
// @Failure      401   {object}  response.Response  "UNAUTHORIZED — invalid username or password"
// @Router       /auth/login [post]
func (this *LoginHandler) Handle(c ctx.Context) error {
	input, err := ctx.GetBody[dto.LoginInput](c)
	if err != nil {
		return err
	}

	output, err := this.uc.Invoke(c.GetContext(), *input)
	if err != nil {
		return err
	}
	return c.Success(http.StatusOK, output)
}
