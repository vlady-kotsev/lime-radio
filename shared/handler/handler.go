package handler

import (
	"github.com/vlady-kotsev/lime-radio/shared/server"
)

func RegisterHandlers(handlers []Handlerer, server server.Serverer) {
	app := server.GetApp()
	for _, handler := range handlers {
		handler.RegisterRoute(app)
	}
}
