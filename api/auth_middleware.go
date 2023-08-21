package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	log.Println("this is Atuh MiddleWare")
	session_token := c.Cookies("session_token")
	if session_token == "" {
		return ErrUnAuthorized()
	}
	return c.Next()
}
