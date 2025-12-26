package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
)

const (
	CookieKey  string = "auth_token"
	QueryKey   string = "token"
	AuthHeader string = "Authorization"
	AuthPrefix string = "Bearer"
)

type AuthMiddleware struct {
	jwtService auth.JWTServicer
}

func NewAuthMiddleware(jwtService auth.JWTServicer) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (am *AuthMiddleware) ImposeAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "OPTIONS" {
			return c.Next()
		}

		if c.Path() == "/token" {
			return c.Next()
		}

		var tokenString string

		tokenString = c.Cookies(CookieKey)
		if tokenString == "" {
			authHeader := c.Get(AuthHeader)
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != AuthPrefix {
					return c.Status(401).JSON(fiber.Map{
						"error": "Invalid authorization header format",
					})
				}
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = c.Query(QueryKey)
		}

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Missing authorization token",
			})
		}

		claims, err := am.jwtService.ValidateToken(tokenString)
		if err != nil {
			if c.Cookies(CookieKey) != "" {
				c.Cookie(&fiber.Cookie{
					Name:     CookieKey,
					Value:    "",
					Expires:  time.Now().Add(-time.Hour),
					HTTPOnly: true,
					Secure:   false,
					SameSite: "Lax",
				})
			}

			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
