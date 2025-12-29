package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/vlady-kotsev/lime-radio/dj/internal/middleware"
	"github.com/vlady-kotsev/lime-radio/dj/internal/service/payment"
	"github.com/vlady-kotsev/lime-radio/dj/internal/service/transaction"
	"github.com/vlady-kotsev/lime-radio/shared/config"
	sharedmiddleware "github.com/vlady-kotsev/lime-radio/shared/middleware"
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

func NewServer(lc fx.Lifecycle, logger *zap.Logger, auth auth.JWTServicer, config config.Configer, paymentService payment.PaymentServicer, transactionService transaction.TransactionServicer) *Server {
	app := fiber.New()

	allowedOrigins := "*"
	allowCredentials := false
	if len(config.GetAllowedOrigins()) > 0 {
		allowedOrigins = strings.Join(config.GetAllowedOrigins(), ",")
		allowCredentials = true
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Payment",
		AllowCredentials: allowCredentials,
		ExposeHeaders:    "Content-Length,Content-Range,Payment-Required",
		MaxAge:           86400,
	}))

	app.Use(sharedmiddleware.NewAuthMiddleware(auth).ImposeAuth())
	app.Use(middleware.NewPaymentMiddleware(paymentService, transactionService).ImposePayment())
	s := &Server{
		fiber:  app,
		logger: logger,
		port:   fmt.Sprintf(":%d", config.GetPort()),
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
