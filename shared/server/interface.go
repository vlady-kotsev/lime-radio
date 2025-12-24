package server

import "github.com/gofiber/fiber/v2"

type Serverer interface {
	GetApp() *fiber.App
}
