package main

import (
	"log"

	"chatemotes/internal/api"
	"chatemotes/internal/database"
	"chatemotes/internal/resourcepack"
	"chatemotes/internal/services"
)

func run() error {
	database := database.New()
	resourcepack := resourcepack.New(database)

	services := &services.Services{
		ResoucePack: resourcepack,
	}

	app := api.New(services)

	return app.Listen(":3000")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
