package main

import (
	"log"
	"math/rand"
	"time"

	"chatemotes/internal/api"
	"chatemotes/internal/database"
	"chatemotes/internal/logic"
)

var _packFilename = "./pack/resourcepack.zip"

func run() error {
	rand.Seed(time.Now().UnixNano())
	database := database.New()
	logic := logic.New(_packFilename, database)
	return api.New(logic).Listen(":3000")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
