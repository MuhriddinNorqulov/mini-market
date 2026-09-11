package auth

import (
	"net/http"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/usecases/authusecases"
	"mini-market/src/core/domain/ports/httpport/ctx"
)

type RegisterHandler struct {
	uc *authusecases.RegisterUseCase
}

// @inject
func NewRegisterHandler(uc *authusecases.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{uc: uc}
}

// Handle godoc
// @Tags         Auth
// @Summary      Register
// @Description  Creates a new user with a username and password and returns an access token.
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterInput  true  "Username, password, and password confirmation"
// @Success      201   {object}  response.Response{payload=dto.AuthOutput}
// @Failure      400   {object}  response.Response  "BAD_REQUEST — validation error, e.g. password and confirm_password don't match"
// @Failure      409   {object}  response.Response  "CONFLICT — username is already taken"
// @Router       /auth/register [post]
func (this *RegisterHandler) Handle(c ctx.Context) error {
	input, err := ctx.GetBody[dto.RegisterInput](c)
	if err != nil {
		return err
	}

	output, err := this.uc.Invoke(c.GetContext(), *input)
	if err != nil {
		return err
	}
	return c.Success(http.StatusCreated, output)
}
