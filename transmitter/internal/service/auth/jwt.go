package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
)

type JWTService struct {
	secret []byte
}

var _ JWTServicer = (*JWTService)(nil)

func NewJWTService(config *config.Config) (*JWTService, error) {
	if config.Auth.SharedSecret.IsEmpty() {
		return nil, fmt.Errorf("jwt secret not configured")
	}
	return &JWTService{
		secret: config.Auth.SharedSecret.Bytes(),
	}, nil
}

// GenerateToken creates a new JWT token with specified expiration
func (j *JWTService) GenerateToken(expiration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "lime-radio",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken checks if a JWT token is valid and not expired
func (j *JWTService) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
