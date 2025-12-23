package main

import (
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
	songrepository "github.com/vlady-kotsev/lime-radio/transmitter/internal/repository/song"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/server"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/auth"
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
			server.NewServer,
			config.Load,
			fx.Annotate(repository.NewStorage, fx.As(new(repository.Storager))),
			fx.Annotate(radio.NewRadio, fx.As(new(radio.Radioer))),
			fx.Annotate(songrepository.NewSongRepository, fx.As(new(songrepository.SongRepositorer))),
			fx.Annotate(radio.NewPlaylist, fx.As(new(radio.Playlister))),
			fx.Annotate(auth.NewJWTService, fx.As(new(auth.JWTServicer))),
		),
		handler.ProvideHandlers(),
		fx.Invoke(
			fx.Annotate(
				handler.RegisterHandlers,
				fx.ParamTags(`group:"handlers"`),
			),
		),
	).Run()
}
