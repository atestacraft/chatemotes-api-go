package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"chatemotes/internal/logic"
)

func setEmoteRoutes(app fiber.Router, logic logic.Logic) {
	app.Get("/emotes", func(c *fiber.Ctx) error {
		emotes := logic.GetEmotes()
		return c.JSON(emotes)
	})

	app.Post("/emotes", func(c *fiber.Ctx) error {
		var body struct {
			Url  string `json:"url" xml:"url" form:"url"`
			Name string `json:"name" xml:"name" form:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err,
				})
		}

		response, err := logic.AddEmote(body.Url, body.Name)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"error": err,
				})
		}

		return c.JSON(response)
	})

	app.Patch("/emotes/:name", func(c *fiber.Ctx) error {
		var body struct {
			Name string `json:"name"` // new emote name
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err,
				})
		}

		emoteName := c.Params("name")
		_, err := logic.GetEmoteByName(emoteName)
		if err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(fiber.Map{
					"error": err,
				})
		}

		logic.UpdateEmote(emoteName, body.Name)
		return c.JSON(fiber.Map{"success": true})
	})

	app.Get("/emotes/:name", func(c *fiber.Ctx) error {
		emote, err := logic.GetEmoteByName(c.Params("name"))
		if err != nil {
			return c.Status(http.StatusNotFound).
				JSON(fiber.Map{
					"error": err,
				})
		}

		return c.JSON(emote)
	})

	app.Delete("/emotes/:name", func(c *fiber.Ctx) error {
		err := logic.RemoveEmoteByName(c.Params("name"))
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(fiber.Map{
					"error": err,
				})
		}
		return c.JSON(fiber.Map{
			"success": true,
		})
	})
}

func New(logic logic.Logic) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/hash", func(c *fiber.Ctx) error {
		type response struct {
			Hash string `json:"hash"`
		}
		return c.JSON(response{
			Hash: logic.GetHash(),
		})
	})
	setEmoteRoutes(app, logic)
	app.Get("/pack", func(c *fiber.Ctx) error {
		return c.Download(logic.GetPackFilename())
	})

	return app
}
