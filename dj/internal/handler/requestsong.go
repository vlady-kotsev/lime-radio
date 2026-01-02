package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/vlady-kotsev/lime-radio/dj/internal/domain"
	"github.com/vlady-kotsev/lime-radio/dj/internal/events"
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
	publisher  events.EventPublisherer
}

var _ handler.Handlerer = (*RequestSongHandler)(nil)

func NewRequestSongHandler(logger *zap.Logger, jwtService auth.JWTServicer, config config.Configer, publisher events.EventPublisherer) *RequestSongHandler {
	return &RequestSongHandler{
		logger:     logger,
		jwtService: jwtService,
		config:     config,
		publisher:  publisher,
		path:       "/request",
	}
}

func (gt *RequestSongHandler) Handle(c *fiber.Ctx) error {
	var songRequest domain.SongRequest
	err := c.BodyParser(&songRequest)
	if err != nil {
		return err
	}

	err = gt.publisher.PublishMessage(songRequest.ID.String())
	if err != nil {
		return err
	}

	gt.logger.Info("Message published")
	return nil
}

func (h *RequestSongHandler) RegisterRoute(app *fiber.App) {
	app.Post(h.path, h.Handle)
}
