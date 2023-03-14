package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"chatemotes/internal/api/emotes"
	"chatemotes/internal/api/hash"
	"chatemotes/internal/api/pack"
	"chatemotes/internal/database"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"
)

func run() error {
	database := database.New()
	resourcepack := resourcepack.New(database)

	services := &services.Services{
		ResoucePack: resourcepack,
		Database:    database,
	}

	app := fiber.New()
	app.Use(logger.New())

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
