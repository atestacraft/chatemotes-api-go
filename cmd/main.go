package main

import (
	"log"
	"math/rand"
	"time"

	"chatemotes/internal/api"
	"chatemotes/internal/database"
	"chatemotes/internal/logic"

	"github.com/joho/godotenv"
	"github.com/rprtr258/xerr"
)

func run() error {
	rand.Seed(time.Now().UnixNano())

	env, err := godotenv.Read(".env")
	if err != nil {
		return err
	}

	resourcePackPath := env["RESOURCEPACK_PATH"]
	if resourcePackPath == "" {
		return xerr.NewM("RESOURCEPACK_PATH is not set")
	}

	database := database.New()
	logic := logic.New(resourcePackPath, database)
	return api.New(logic).Listen(env["SERVER_ADDR"])
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
