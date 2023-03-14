package emotes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"chatemotes/internal/logic"
)

type EmoteBody struct {
	Url  string `json:"url" xml:"url" form:"url"`
	Name string `json:"name" xml:"name" form:"name"`
}

func NewRoute(app fiber.Router, logic logic.Logic) {
	app.Get("/emotes", func(c *fiber.Ctx) error {
		res, err := logic.GetEmotes()
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		return c.JSON(res)
	})

	app.Post("/emotes", func(c *fiber.Ctx) error {
		var emoteBody EmoteBody
		if err := c.BodyParser(&emoteBody); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		response, err := logic.AddEmote(emoteBody.Url, emoteBody.Name)
		if err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		return c.JSON(response)
	})

	// TODO: :name ???
	app.Patch("/emotes", func(c *fiber.Ctx) error {
		var body struct {
			Url  string `json:"url"`
			Name string `json:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		response, err := logic.UpdateEmote(body.Name)
		if err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		return c.JSON(response)
	})

	app.Get("/emotes/:name", func(c *fiber.Ctx) error {
		emote, err := logic.GetEmoteByName(c.Params("name"))
		if err != nil {
			return c.Status(http.StatusNotFound).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		if emote == nil {
			return c.SendStatus(http.StatusNotFound)
		}

		return c.JSON(emote)
	})

	app.Delete("/emotes/:name", func(c *fiber.Ctx) error {
		err := logic.RemoveEmoteByName(c.Params("name"))
		if err != nil {
			return c.Status(http.StatusNotFound).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}
		return c.JSON(fiber.Map{
			"success": true,
		})
	})
}
