package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type UpdateSongsHandler struct {
	pl     radio.PlaylistServicer
	logger *zap.Logger
	path   string
}

var _ handler.Handlerer = (*UpdateSongsHandler)(nil)

func NewUpdateSongsHandler(pl radio.PlaylistServicer, logger *zap.Logger) *UpdateSongsHandler {
	return &UpdateSongsHandler{
		pl:     pl,
		logger: logger,
		path:   "/refresh",
	}
}

func (h *UpdateSongsHandler) Handle(c *fiber.Ctx) error {
	err := h.pl.UpdateSongs()
	if err != nil {
		return err
	}
	songs, err := h.pl.GetAllSongs()
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
