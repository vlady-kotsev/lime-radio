package main

import (
	"github.com/vlady-kotsev/lime-radio/dj/internal/config"
	"github.com/vlady-kotsev/lime-radio/dj/internal/handler"
	"github.com/vlady-kotsev/lime-radio/dj/internal/server"
	sharedconfig "github.com/vlady-kotsev/lime-radio/shared/config"
	sharedserver "github.com/vlady-kotsev/lime-radio/shared/server"

	sharedhandler "github.com/vlady-kotsev/lime-radio/shared/handler"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(
			zap.NewProduction,
			fx.Annotate(server.NewServer, fx.As(new(sharedserver.Serverer))),
			fx.Annotate(config.Load, fx.As(new(sharedconfig.AuthConfiger)), fx.As(fx.Self())),
			fx.Annotate(
				func(cfg *config.Config) (auth.JWTServicer, error) {
					return auth.NewJWTService(cfg.Auth.SharedSecret)
				},
				fx.As(new(auth.JWTServicer)),
			),
		),
		handler.ProvideHandlers(),
		fx.Invoke(
			fx.Annotate(
				sharedhandler.RegisterHandlers,
				fx.ParamTags(`group:"handlers"`),
			),
		),
	).Run()
}
