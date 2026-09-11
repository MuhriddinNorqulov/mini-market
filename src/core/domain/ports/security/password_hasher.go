package security

type PasswordHasher interface {
	Hash(password string) (string, error)
	Equal(hashedPassword, password string) bool
}
