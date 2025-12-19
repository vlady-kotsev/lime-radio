package handler

import (
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/server"
	"go.uber.org/fx"
)

func ProvideHandlers() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				NewStreamHandler,
				fx.As(new(Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				NewGetAllSongsHandler,
				fx.As(new(Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				NewUpdateSongsHandler,
				fx.As(new(Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
		),
	)
}

func RegisterHandlers(handlers []Handlerer, server *server.Server) {
	app := server.GetApp()
	for _, handler := range handlers {
		handler.RegisterRoute(app)
	}
}
