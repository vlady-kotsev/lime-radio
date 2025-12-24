package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/vlady-kotsev/lime-radio/dj/internal/config"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
	"go.uber.org/zap"
)

type RequestSongHandler struct {
	logger     *zap.Logger
	jwtService auth.JWTServicer
	config     *config.Config
	path       string
}

var _ handler.Handlerer = (*RequestSongHandler)(nil)

func NewRequestSongHandler(logger *zap.Logger, jwtService auth.JWTServicer, config *config.Config) *RequestSongHandler {
	return &RequestSongHandler{
		logger:     logger,
		jwtService: jwtService,
		config:     config,
		path:       "/request",
	}
}

func (gt *RequestSongHandler) Handle(c *fiber.Ctx) error {
	c.Status(402)
	return nil
}

func (h *RequestSongHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
