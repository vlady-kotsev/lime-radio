package server

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
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

func NewServer(lc fx.Lifecycle, logger *zap.Logger, storage *repository.Storage, config *config.Config) *Server {
	app := fiber.New()
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
		ExposeHeaders: "Content-Length,Content-Range",
		MaxAge: 86400,
	}))
	
	s := &Server{
		fiber:   app,
		logger:  logger,
		port:    fmt.Sprintf(":%d", config.App.Port),
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
