package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/api/pack"
	"chatemotes/internal/database"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"
)

func run() error {
	app := fiber.New()
	database := database.New()
	resourcepack := resourcepack.New(database)

	services := &services.Services{
		ResoucePack: resourcepack,
		Database:    database,
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
