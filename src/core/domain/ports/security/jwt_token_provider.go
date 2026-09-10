package security

type JwtTokenProvider interface {
	Encode(token *JwtToken) string
	Decode(token string) (*JwtToken, error)
}
