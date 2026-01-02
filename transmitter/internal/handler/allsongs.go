package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

const SearchQueryKey string = "search"

type GetAllSongsHandler struct {
	pl     radio.PlaylistServicer
	logger *zap.Logger
	path   string
}

var _ handler.Handlerer = (*GetAllSongsHandler)(nil)

func NewGetAllSongsHandler(pl radio.PlaylistServicer, logger *zap.Logger) *GetAllSongsHandler {
	return &GetAllSongsHandler{
		pl:     pl,
		logger: logger,
		path:   "/songs",
	}
}

func (h *GetAllSongsHandler) Handle(c *fiber.Ctx) error {
	searchQuery := c.Query(SearchQueryKey)
	var songs []*domain.Song
	var err error
	if searchQuery == "" {
		songs, err = h.pl.GetAllSongs()
	} else {
		songs, err = h.pl.GetAllSongsByTitleOrArtist(searchQuery)
	}
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
