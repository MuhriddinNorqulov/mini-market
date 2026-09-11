package authusecases

import (
	"context"
	"time"

	"mini-market/src/core/application/dto"
	"mini-market/src/core/application/response"
	"mini-market/src/core/application/services"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/security"
	"mini-market/src/core/domain/ports/unitofwork"
)

type RegisterUseCase struct {
	userRepo    repository.UserRepository
	hasher      security.PasswordHasher
	authService *services.UserAuthTokenService
	atomic      unitofwork.Atomic
}

// @inject
func NewRegisterUseCase(
	userRepo repository.UserRepository,
	hasher security.PasswordHasher,
	authService *services.UserAuthTokenService,
	atomic unitofwork.Atomic,
) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo, hasher: hasher, authService: authService, atomic: atomic}
}

func (this *RegisterUseCase) Invoke(ctx context.Context, input dto.RegisterInput) (*dto.AuthOutput, error) {
	if _, err := this.userRepo.GetByUsername(ctx, input.Username); err == nil {
		return nil, response.NewResponse(response.CodeConflict, false, nil, "username is already taken")
	}

	passwordHash, err := this.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	var user *entity.UserEntity
	err = this.atomic.Transaction(func(tx unitofwork.Tx) error {
		created, err := this.userRepo.CreateTx(ctx, tx, &entity.UserEntity{
			Username:  &input.Username,
			FirstName: input.Username,
			Role:      enum.RoleUser,
		})
		if err != nil {
			return err
		}

		if err := this.userRepo.UpdatePasswordTx(ctx, tx, created.ID, passwordHash, time.Now()); err != nil {
			return err
		}

		user = created
		return nil
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := this.authService.GenerateAccessToken(user, "")
	if err != nil {
		return nil, err
	}
	return &dto.AuthOutput{AccessToken: accessToken}, nil
}
