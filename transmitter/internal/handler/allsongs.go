package handler

import (
	"github.com/gofiber/fiber/v2"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type GetAllSongsHandler struct {
	radio  radio.Radioer
	logger *zap.Logger
	path   string
}

var _ Handlerer = (*GetAllSongsHandler)(nil)

func NewGetAllSongsHandler(radio radio.Radioer, logger *zap.Logger) *GetAllSongsHandler {
	return &GetAllSongsHandler{
		radio:  radio,
		logger: logger,
		path:   "/songs",
	}
}

func (h *GetAllSongsHandler) Handle(c *fiber.Ctx) error {
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

func (h *GetAllSongsHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
