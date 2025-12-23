package handler

import (
	"github.com/gofiber/fiber/v2"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type UpdateSongsHandler struct {
	radio  radio.Radioer
	logger *zap.Logger
	path   string
}

var _ Handlerer = (*UpdateSongsHandler)(nil)

func NewUpdateSongsHandler(radio radio.Radioer, logger *zap.Logger) *UpdateSongsHandler {
	return &UpdateSongsHandler{
		radio:  radio,
		logger: logger,
		path:   "/refresh",
	}
}

func (h *UpdateSongsHandler) Handle(c *fiber.Ctx) error {
	err := h.radio.UpdateSongs()
	if err != nil {
		return err
	}
	songs, err := h.radio.GetAllSongs()
	if err != nil {
		h.logger.Error("Failed to get songs", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve songs",
		})
	}

	viewModels := songviewmodel.ToSongViewModels(songs)
	return c.JSON(viewModels)
}

func (h *UpdateSongsHandler) RegisterRoute(app *fiber.App) {
	app.Post(h.path, h.Handle)
}
