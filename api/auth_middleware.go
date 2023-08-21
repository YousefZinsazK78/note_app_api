package api

import (
	"github.com/gofiber/fiber/v2"
)

func (a *Api) AuthMiddleware(c *fiber.Ctx) error {
	session_token := c.Cookies("session_token")
	if session_token == "" {
		return ErrUnAuthorized()
	}
	session, err := a.SessionStorer.GetSession(session_token)
	if err != nil {
		return NewError(fiber.StatusBadRequest, err.Error())
	}
	if session.IsExpired() {
		return ErrUnAuthorized()
	}
	return c.Next()
}
