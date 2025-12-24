package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/vlady-kotsev/lime-radio/dj/internal/config"
	"github.com/vlady-kotsev/lime-radio/shared/middleware"
	"github.com/vlady-kotsev/lime-radio/shared/server"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	fiber  *fiber.App
	logger *zap.Logger
	port   string
}

var _ server.Serverer = (*Server)(nil)

func NewServer(lc fx.Lifecycle, logger *zap.Logger, auth auth.JWTServicer, config *config.Config) *Server {
	app := fiber.New()

	allowedOrigins := "*"
	if len(config.Auth.AllowedOrigins) > 0 {
		allowedOrigins = strings.Join(config.Auth.AllowedOrigins, ",")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Range",
		MaxAge:           86400,
	}))

	app.Use(middleware.JWTAuth(auth))

	s := &Server{
		fiber:  app,
		logger: logger,
		port:   fmt.Sprintf(":%d", config.App.Port),
	}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Starting HTTP server at", zap.String("addr", s.port))
				go func() {
					if err := s.fiber.Listen(s.port); err != nil {
						logger.Fatal("Failed to start server", zap.Error(err))
					}
				}()
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
