package emotes

import (
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
)

func NewRoute(app fiber.Router, services *services.Services) {
	app.Get("/emotes", func(c *fiber.Ctx) error {
		return c.JSON(services.ResoucePack.GetEmotes())
	})

	app.Get("/emotes/:name", func(c *fiber.Ctx) error {
		response, err := services.ResoucePack.GetEmoteByName(c.Params("name"))
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(response)
	})
}
