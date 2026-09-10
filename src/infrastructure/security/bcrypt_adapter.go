package security

import (
	"mini-market/src/core/domain/ports/security"

	"golang.org/x/crypto/bcrypt"
)

type BcryptAdapter struct {
	cost int
}

// @inject
func NewBcryptAdapter() security.PasswordHasher {
	return &BcryptAdapter{cost: bcrypt.DefaultCost}
}

func (b *BcryptAdapter) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptAdapter) Equal(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
