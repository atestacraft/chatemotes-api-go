package main

import (
	"chatemotes/internal/api/hash"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	services := &services.Services{
		ResoucePack: resourcepack.New(),
	}

	hash.NewRoute(app, services)

	app.Listen(":3000")
}
