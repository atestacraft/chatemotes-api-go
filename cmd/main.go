package main

import (
	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/api/pack"
	"chatemotes/internal/database"
	emote_resolver "chatemotes/internal/emote"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	database := database.New()
	resourcepack := resourcepack.New(database)
	emoteResolver := emote_resolver.New()

	services := &services.Services{
		ResoucePack:   resourcepack,
		Database:      database,
		EmoteResolver: emoteResolver,
	}

	hash.NewRoute(app, services)
	emotes.NewRoute(app, services)
	pack.NewRoute(app)

	app.Listen(":3000")
}
