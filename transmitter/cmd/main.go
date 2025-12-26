package main

import (
	sharedconfig "github.com/vlady-kotsev/lime-radio/shared/config"
	sharedhandler "github.com/vlady-kotsev/lime-radio/shared/handler"
	sharedserver "github.com/vlady-kotsev/lime-radio/shared/server"
	"github.com/vlady-kotsev/lime-radio/shared/service/auth"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
	songrepository "github.com/vlady-kotsev/lime-radio/transmitter/internal/repository/song"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/server"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
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
			fx.Annotate(config.Load, fx.As(new(sharedconfig.AuthConfiger)), fx.As(new(sharedconfig.Configer)), fx.As(new(config.RadioConfiger))),
			fx.Annotate(repository.NewStorage, fx.As(new(repository.Storager))),
			fx.Annotate(radio.NewRadio, fx.As(new(radio.RadioServicer))),
			fx.Annotate(songrepository.NewSongRepository, fx.As(new(songrepository.SongRepositorer))),
			fx.Annotate(radio.NewPlaylist, fx.As(new(radio.PlaylistServicer))),
			fx.Annotate(auth.NewJWTService, fx.As(new(auth.JWTServicer))),
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
