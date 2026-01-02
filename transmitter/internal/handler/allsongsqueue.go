package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type GetAllSongsInQueueHandler struct {
	pl     radio.PlaylistServicer
	logger *zap.Logger
	path   string
}

var _ handler.Handlerer = (*GetAllSongsInQueueHandler)(nil)

func NewGetAllSongsInQueueHandler(pl radio.PlaylistServicer, logger *zap.Logger) *GetAllSongsInQueueHandler {
	return &GetAllSongsInQueueHandler{
		pl:     pl,
		logger: logger,
		path:   "/queue",
	}
}

func (h *GetAllSongsInQueueHandler) Handle(c *fiber.Ctx) error {
	songs := h.pl.GetAllSongsInQueue()

	viewModels := songviewmodel.ToSongViewModels(songs)
	return c.JSON(viewModels)
}

func (h *GetAllSongsInQueueHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
