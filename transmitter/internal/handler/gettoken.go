package handler

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/auth"
	"go.uber.org/zap"
)

type GetTokenHandler struct {
	logger     *zap.Logger
	jwtService *auth.JWTService
	config     *config.Config
	path       string
}

var _ Handlerer = (*GetTokenHandler)(nil)

func NewGetTokenHandler(logger *zap.Logger, jwtService *auth.JWTService, config *config.Config) *GetTokenHandler {
	return &GetTokenHandler{
		logger:     logger,
		jwtService: jwtService,
		config:     config,
		path:       "/token",
	}
}

// GetAuthToken validates origin and returns a short-lived JWT token
func (gt *GetTokenHandler) Handle(c *fiber.Ctx) error {
	origin := c.Get("Origin")
	if origin == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Missing Origin header",
		})
	}

	if !slices.Contains(gt.config.Auth.AllowedOrigins, origin) {
		gt.logger.Warn("Unauthorized origin", zap.String("origin", origin))
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized origin",
		})

	}

	expiration := time.Duration(gt.config.Auth.TokenExpirationMinutes) * time.Minute
	token, err := gt.jwtService.GenerateToken(expiration)
	if err != nil {
		gt.logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	expiresAt := time.Now().Add(expiration)
	gt.logger.Info("Token generated", zap.String("origin", origin), zap.Time("expires", expiresAt))

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"expires": expiresAt.Unix(),
		"token":   token,
	})
}

func (h *GetTokenHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
