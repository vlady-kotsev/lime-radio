package handler

import "github.com/gofiber/fiber/v2"

type Handlerer interface {
	Handle(c *fiber.Ctx) error
	RegisterRoute(app *fiber.App)
}
