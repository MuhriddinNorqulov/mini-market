package security

import "time"

type JwtToken struct {
	Exp     time.Time
	Subject string
	Payload []byte
}
