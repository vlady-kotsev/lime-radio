package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
	songviewmodel "github.com/vlady-kotsev/lime-radio/transmitter/internal/handler/viewmodel"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

const (
	SearchQueryKey       string = "search"
	EmptyPaginationValue string = "0"
)

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
	page, err := strconv.Atoi(c.Query("page", EmptyPaginationValue))
	if err != nil {
		h.logger.Error("Failed to parse page query param", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse page query param",
		})
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", EmptyPaginationValue))
	if err != nil {
		h.logger.Error("Failed to parse page_size query param", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse page_size query param",
		})
	}

	paginationParams := domain.NewPaginationParams(page, pageSize)

	searchQuery := c.Query(SearchQueryKey)
	if searchQuery == "" {
		result, err := h.pl.GetSongsPaginated(paginationParams)
		if err != nil {
			h.logger.Error("Failed to get paginated songs", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve songs",
			})
		}

		return c.JSON(songviewmodel.NewPaginatedSongsViewModel(result.Data, result.Page, result.HasNext))
	} else {
		result, err := h.pl.GetSongsByTitleOrArtist(searchQuery, paginationParams)
		if err != nil {
			h.logger.Error("Failed to get paginated search results", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve songs",
			})
		}

		return c.JSON(songviewmodel.NewPaginatedSongsViewModel(result.Data, result.Page, result.HasNext))
	}
}

func (h *GetAllSongsHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
}
