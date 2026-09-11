package authusecases

import (
	"context"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/response"
	"mini-market/src/core/application/services"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/security"
)

type LoginUseCase struct {
	userRepo    repository.UserRepository
	hasher      security.PasswordHasher
	authService *services.UserAuthTokenService
}

// @inject
func NewLoginUseCase(
	userRepo repository.UserRepository,
	hasher security.PasswordHasher,
	authService *services.UserAuthTokenService,
) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, hasher: hasher, authService: authService}
}

func (this *LoginUseCase) Invoke(ctx context.Context, input dto.LoginInput) (*dto.AuthOutput, error) {
	user, passwordHash, err := this.userRepo.GetByUsernameWithPassword(ctx, input.Username)
	if err != nil {
		return nil, response.NewResponse(response.CodeUnauthorized, false, nil, "invalid username or password")
	}

	if passwordHash == nil || !this.hasher.Equal(*passwordHash, input.Password) {
		return nil, response.NewResponse(response.CodeUnauthorized, false, nil, "invalid username or password")
	}

	accessToken, err := this.authService.GenerateAccessToken(user, "")
	if err != nil {
		return nil, err
	}
	return &dto.AuthOutput{AccessToken: accessToken}, nil
}
