package authService

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
)

func (a authService) HandlerLogout(c *fiber.Ctx) error {
	accessToken := c.Locals("token").(*entity.JsonWebToken)

	if accessToken.ID == "" {
		panic(errors.New("logout failed: no access token found"))
	}

	ref := accessToken.Ref
	err := a.jsonWebTokenRepository.RevokeToken(ref)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
