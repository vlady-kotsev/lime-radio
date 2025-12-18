package handler

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/audio"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type StreamHandler struct {
	station *radio.Station
	logger  *zap.Logger
	path    string
}

func NewStreamHandler(station *radio.Station, logger *zap.Logger) *StreamHandler {
	return &StreamHandler{
		station: station,
		logger:  logger,
		path:    "/stream",
	}
}

func (h *StreamHandler) HandleStream(c *fiber.Ctx) error {
	c.Set("Content-Type", "audio/wav")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		client := h.station.AddClient()
		defer h.station.RemoveClient(client)

		// Send WAV header first
		if h.station.GetSampleRate() > 0 {
			header := audio.CreateWAVHeader(h.station.GetSampleRate())
			if _, err := w.Write(header); err != nil {
				h.logger.Error("Error writing WAV header", zap.Error(err))
				return
			}
			w.Flush()
		}

		// Stream audio data
		for data := range client {
			if _, err := w.Write(data); err != nil {
				h.logger.Debug("Client disconnected", zap.Error(err))
				return
			}
			w.Flush()
		}
	})

	return nil
}

func (h *StreamHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.HandleStream)
}
