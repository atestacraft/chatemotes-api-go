package emotes

import (
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
)

type EmoteBody struct {
	Url  string `json:"url" xml:"url" form:"url"`
	Name string `json:"name" xml:"name" form:"name"`
}

func NewRoute(app fiber.Router, services *services.Services) {
	app.Get("/emotes", func(c *fiber.Ctx) error {
		return c.JSON(services.ResoucePack.GetEmotes())
	})

	app.Post("/emotes", func(c *fiber.Ctx) error {
		emoteBody := new(EmoteBody)

		if err := c.BodyParser(emoteBody); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		response, err := services.ResoucePack.AddEmote(emoteBody.Url, emoteBody.Name)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(response)
	})

	app.Patch("/emotes", func(c *fiber.Ctx) error {
		var body struct {
			Url  string `json:"url"`
			Name string `json:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		response, err := services.ResoucePack.UpdateEmote(body.Url, body.Name)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(response)
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

	app.Delete("/emotes/:name", func(c *fiber.Ctx) error {
		err := services.ResoucePack.RemoveEmoteByName(c.Params("name"))
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
		})
	})
}
