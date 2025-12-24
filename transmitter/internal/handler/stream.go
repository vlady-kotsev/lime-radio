package handler

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type StreamHandler struct {
	radio  radio.Radioer
	logger *zap.Logger
	path   string
}

var _ handler.Handlerer = (*StreamHandler)(nil)

func NewStreamHandler(radio radio.Radioer, logger *zap.Logger) *StreamHandler {
	return &StreamHandler{
		radio:  radio,
		logger: logger,
		path:   "/stream",
	}
}

func (h *StreamHandler) Handle(c *fiber.Ctx) error {
	c.Set("Content-Type", "audio/wav")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	c.Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		client := h.radio.AddClient()
		defer h.radio.RemoveClient(client)

		if h.radio.GetSampleRate() > 0 {
			header, err := radio.CreateWAVHeader(h.radio.GetSampleRate())
			if err != nil {
				h.logger.Error("Error creating WAV header", zap.Error(err))
				return
			}
			if _, err = w.Write(header); err != nil {
				h.logger.Error("Error writing WAV header", zap.Error(err))
				return
			}
			err = w.Flush()
			if err != nil {
				h.logger.Error("Error flushing WAV header", zap.Error(err))
				return
			}
		}

		for data := range client {
			if _, err := w.Write(data); err != nil {
				h.logger.Debug("Client disconnected", zap.Error(err))
				return
			}
			err := w.Flush()
			if err != nil {
				h.logger.Error("Error flushing data", zap.Error(err))
				return
			}
		}
	})

	return nil
}

func (h *StreamHandler) RegisterRoute(app *fiber.App) {
	app.Get(h.path, h.Handle)
	app.Options(h.path, func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Max-Age", "86400")
		return c.SendStatus(204)
	})
}
