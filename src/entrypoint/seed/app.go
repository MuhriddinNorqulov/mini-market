package seed

import (
	"context"
	"fmt"
	"time"

	"mini-market/src/core/application/services"
	"mini-market/src/core/domain/entity"
	"mini-market/src/core/domain/entity/enum"
	"mini-market/src/core/domain/ports/repository"
	"mini-market/src/core/domain/ports/security"
	"mini-market/src/core/domain/ports/unitofwork"
)

type App struct {
	userRepository repository.UserRepository
	authService    *services.UserAuthTokenService
	atomic         unitofwork.Atomic
	hasher         security.PasswordHasher
}

// @inject
func NewApp(
	userRepository repository.UserRepository,
	authService *services.UserAuthTokenService,
	atomic unitofwork.Atomic,
	hasher security.PasswordHasher,
) *App {
	return &App{userRepository: userRepository, authService: authService, atomic: atomic, hasher: hasher}
}

func (this *App) IssueTestToken(username string) (string, error) {
	var user *entity.UserEntity

	err := this.atomic.Transaction(func(tx unitofwork.Tx) error {
		created, err := this.userRepository.CreateTx(context.Background(), tx, &entity.UserEntity{
			Username:  &username,
			FirstName: "Concurrency",
			Role:      enum.RoleUser,
		})
		if err != nil {
			return err
		}
		user = created
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("create test user: %w", err)
	}

	token, err := this.authService.GenerateAccessToken(user, "")
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	return token, nil
}

type DefaultUserSpec struct {
	Username string
	Password string
	Role     enum.Role
}

func (this *App) SeedDefaultUsers(ctx context.Context, specs []DefaultUserSpec) error {
	for _, spec := range specs {
		if _, err := this.userRepository.GetByUsername(ctx, spec.Username); err == nil {
			continue
		}

		passwordHash, err := this.hasher.Hash(spec.Password)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", spec.Username, err)
		}

		username := spec.Username
		err = this.atomic.Transaction(func(tx unitofwork.Tx) error {
			created, err := this.userRepository.CreateTx(ctx, tx, &entity.UserEntity{
				Username:  &username,
				FirstName: spec.Username,
				Role:      spec.Role,
			})
			if err != nil {
				return err
			}
			return this.userRepository.UpdatePasswordTx(ctx, tx, created.ID, passwordHash, time.Now())
		})
		if err != nil {
			return fmt.Errorf("create default user %s: %w", spec.Username, err)
		}
	}
	return nil
}
