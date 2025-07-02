package authService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
)

func (a authService) HandlerGetUserInfo(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.Users)
	if user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "User context missing",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user.ToResponse(a.encryptorRepository.Decrypt))
}
