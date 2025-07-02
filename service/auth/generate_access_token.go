package authService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
)

func (a authService) HandlerRefreshAccessToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Users)
	refreshToken := c.Locals("token").(*entity.JsonWebToken)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	token, err := a.jsonWebTokenRepository.GenerateAccessToken(user, refreshToken.ID)
	if err != nil {
		panic(err)
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
