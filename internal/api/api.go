package api

import (
	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func New(services *services.Services) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	hash.NewRoute(app, services)
	emotes.NewRoute(app, services)
	app.Static("/pack", "./pack", fiber.Static{
		Index:          "resourcepack.zip",
		Download:       true,
		Compress:       false,
		ByteRange:      false,
		Browse:         false,
		CacheDuration:  0,
		MaxAge:         0,
		ModifyResponse: nil,
		Next:           nil,
	})

	return app
}
