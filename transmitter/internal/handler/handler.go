package handler

import (
	"github.com/vlady-kotsev/lime-radio/shared/handler"
	"go.uber.org/fx"
)

func ProvideHandlers() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				NewStreamHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				NewGetAllSongsHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				NewUpdateSongsHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				handler.NewGetTokenHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				NewGetAllSongsInQueueHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
			fx.Annotate(
				NewQueueCountHandler,
				fx.As(new(handler.Handlerer)),
				fx.ResultTags(`group:"handlers"`),
			),
		),
	)
}
