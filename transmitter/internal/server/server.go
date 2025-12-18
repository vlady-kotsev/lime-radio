package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	fiber   *fiber.App
	logger  *zap.Logger
	port    string
	storage *repository.Storage
}

func NewServer(lc fx.Lifecycle, logger *zap.Logger, storage *repository.Storage) *Server {
	s := &Server{
		fiber:   fiber.New(),
		logger:  logger,
		port:    ":8080",
		storage: storage,
	}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting HTTP server at", zap.String("addr", s.port))
				go s.fiber.Listen(s.port)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.fiber.Shutdown()
			},
		})

	return s
}

func (s *Server) GetApp() *fiber.App {
	return s.fiber
}

func (s *Server) Run() {

}
