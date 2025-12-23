package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/auth"
)

// JWTAuth creates middleware that validates JWT tokens
func JWTAuth(jwtService auth.JWTServicer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "OPTIONS" {
			return c.Next()
		}

		if c.Path() == "/token" {
			return c.Next()
		}

		var tokenString string

		tokenString = c.Cookies("auth_token")
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					return c.Status(401).JSON(fiber.Map{
						"error": "Invalid authorization header format",
					})
				}
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Missing authorization token",
			})
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			if c.Cookies("auth_token") != "" {
				c.Cookie(&fiber.Cookie{
					Name:     "auth_token",
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
