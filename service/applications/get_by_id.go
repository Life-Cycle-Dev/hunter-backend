package applicationsService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

func (a applicationsService) HandlerGetApplicationById(c *fiber.Ctx) error {
	permissions := c.Locals("permissions").([]string)
	err := util.CheckAllow(permissions, "application-view")
	if err != nil {
		panic(err)
	}

	id := c.Params("id")
	result, err := a.applicationsRepository.FindById(id)

	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
