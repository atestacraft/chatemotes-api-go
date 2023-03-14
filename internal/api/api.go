package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/logic"
)

func New(logic logic.Logic) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	hash.NewRoute(app, logic)
	emotes.NewRoute(app, logic)
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
