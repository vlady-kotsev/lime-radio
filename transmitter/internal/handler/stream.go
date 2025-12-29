package handler

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/zap"
)

type StreamHandler struct {
	radio  radio.RadioServicer
	logger *zap.Logger
	path   string
}

var _ handler.Handlerer = (*StreamHandler)(nil)

func NewStreamHandler(radio radio.RadioServicer, logger *zap.Logger) *StreamHandler {
	return &StreamHandler{
		radio:  radio,
		logger: logger,
		path:   "/stream",
	}
}

func (h *StreamHandler) Handle(c *fiber.Ctx) error {
	c.Set("Content-Type", "audio/wav")
	c.Set("Cache-Control", "public, max-age=3600")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("Accept-Ranges", "none")
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET")
	c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		connection := h.radio.AddClient()
		defer h.radio.RemoveClient(connection.ID)

		h.logger.Info("Client connected", zap.String("client_id", connection.ID.String()))

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

		for data := range connection.DataChan {

			if _, err := w.Write(data); err != nil {
				h.logger.Debug("Client disconnected",
					zap.String("client_id", connection.ID.String()),
					zap.Error(err))
				return
			}
			err := w.Flush()
			if err != nil {
				h.logger.Error("Error flushing data",
					zap.String("client_id", connection.ID.String()),
					zap.Error(err))
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
