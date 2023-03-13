package main

import (
	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/api/pack"
	"chatemotes/internal/database"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	database := database.New()
	resourcepack := resourcepack.New(database)

	// resourcepack.AddEmote("test")
	// resourcepack.AddEmote("test2")
	// resourcepack.AddEmote("hello")

	services := &services.Services{
		ResoucePack: resourcepack,
		Database:    database,
	}

	hash.NewRoute(app, services)
	emotes.NewRoute(app, services)
	pack.NewRoute(app)

	app.Listen(":3000")
}
