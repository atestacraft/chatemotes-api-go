package hash

import (
	"chatemotes/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Hash string `json:"hash"`
}

func NewRoute(server fiber.Router, services *services.Services) {
	server.Get("/hash", func(c *fiber.Ctx) error {
		response := &Response{
			Hash: services.ResoucePack.GetHash(),
		}

		return c.JSON(response)
	})
}
