package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration, tokenType TokenType) (string, *Payload, error)

	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}