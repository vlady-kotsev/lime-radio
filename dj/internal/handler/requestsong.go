package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/vlady-kotsev/lime-radio/dj/internal/domain"
	"github.com/vlady-kotsev/lime-radio/shared/config"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
	"go.uber.org/zap"
)

type RequestSongHandler struct {
	logger     *zap.Logger
	jwtService auth.JWTServicer
	config     config.Configer
	path       string
}

var _ handler.Handlerer = (*RequestSongHandler)(nil)

func NewRequestSongHandler(logger *zap.Logger, jwtService auth.JWTServicer, config config.Configer) *RequestSongHandler {
	return &RequestSongHandler{
		logger:     logger,
		jwtService: jwtService,
		config:     config,
		path:       "/request",
	}
}

func (gt *RequestSongHandler) Handle(c *fiber.Ctx) error {
	var songRequest domain.SongRequest
	if err := c.BodyParser(&songRequest); err == nil {
		// TODO Send SongRequested event
	}

	return nil
}

func (h *RequestSongHandler) RegisterRoute(app *fiber.App) {
	app.Post(h.path, h.Handle)
}
