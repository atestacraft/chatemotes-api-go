package hash

import (
	"github.com/gofiber/fiber/v2"

	"chatemotes/internal/logic"
)

type Response struct {
	Hash string `json:"hash"`
}

func NewRoute(app fiber.Router, logic logic.Logic) {
	app.Get("/hash", func(c *fiber.Ctx) error {
		response := &Response{
			Hash: logic.GetHash(),
		}

		return c.JSON(response)
	})
}
