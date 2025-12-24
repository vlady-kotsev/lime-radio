package handler

import (
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"go.uber.org/fx"
)

func ProvideHandlers() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				NewRequestSongHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				handler.NewGetTokenHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
		),
	)
}
