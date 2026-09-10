package security

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenClaim struct {
	jwt.RegisteredClaims

	Payload json.RawMessage `json:"payload"`
}
