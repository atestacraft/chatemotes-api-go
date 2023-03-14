package main

import (
	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/api/pack"
	"chatemotes/internal/database"
	emote_resolver "chatemotes/internal/emote"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

func run() error {
	app := fiber.New()
	database := database.New()
	emoteResolver := emote_resolver.New()
	resourcepack := resourcepack.New(database, emoteResolver)

	services := &services.Services{
		ResoucePack:   resourcepack,
		Database:      database,
		EmoteResolver: emoteResolver,
	}

	hash.NewRoute(app, services)
	emotes.NewRoute(app, services)
	pack.NewRoute(app)

	return app.Listen(":3000")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
