package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServicer interface {
	GenerateToken(expiration time.Duration) (string, error)
	ValidateToken(tokenString string) (*jwt.RegisteredClaims, error)
}
