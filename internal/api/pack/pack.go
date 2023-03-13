package pack

import (
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app fiber.Router) {
	app.Static("/pack", "./pack", fiber.Static{
		Index:    "resourcepack.zip",
		Download: true,
	})
}
