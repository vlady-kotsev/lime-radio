package handler

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/vlady-kotsev/lime-radio/shared/config"
	"github.com/vlady-kotsev/lime-radio/shared/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
	"go.uber.org/zap"
)

const CookieKey string = "auth_token"

type GetTokenHandler struct {
	logger     *zap.Logger
	jwtService auth.JWTServicer
	config     config.AuthConfiger
	path       string
}

var _ Handlerer = (*GetTokenHandler)(nil)

func NewGetTokenHandler(logger *zap.Logger, jwtService auth.JWTServicer, config config.AuthConfiger) *GetTokenHandler {
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

	if !slices.Contains(gt.config.GetAllowedOrigins(), origin) {
		gt.logger.Warn("Unauthorized origin", zap.String("origin", origin))
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized origin",
		})

	}

	expiration := time.Duration(gt.config.GetTokenExpirationMinutes()) * time.Minute
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
		Name:     CookieKey,
		Value:    token,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(viewmodel.ToTokenViewModel(true, expiresAt.Unix(), token))
}

func (h *GetTokenHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
