package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type GetQueueCountHandler struct {
	pl     radio.PlaylistServicer
	config config.RadioConfiger
	logger *zap.Logger
	path   string
}

var _ handler.Handlerer = (*GetQueueCountHandler)(nil)

func NewQueueCountHandler(pl radio.PlaylistServicer, logger *zap.Logger, config config.RadioConfiger) *GetQueueCountHandler {
	return &GetQueueCountHandler{
		pl:     pl,
		logger: logger,
		config: config,
		path:   "/queue-count",
	}
}

func (h *GetQueueCountHandler) Handle(c *fiber.Ctx) error {
	songs := h.pl.GetAllSongsInQueue()
	maxQueueCount := h.config.GetQueueMaxCount()
	return c.JSON(songviewmodel.NewQueueCountViewModel(len(songs), maxQueueCount))
}

func (h *GetQueueCountHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
