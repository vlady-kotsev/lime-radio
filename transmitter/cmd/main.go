package main

import (
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/handler"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/server"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/playlist"
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
			radio.NewStation,
			handler.NewStreamHandler,
			repository.NewStorage,
			repository.NewSongRepository,
			playlist.NewPlaylist,
		),
		fx.Invoke(
			func(
				server *server.Server,
				streamHandler *handler.StreamHandler) {
				streamHandler.RegisterRoute(server.GetApp())
			}),
	).Run()
}
