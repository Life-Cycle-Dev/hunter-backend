package applicationsService

import "github.com/gofiber/fiber/v2"

func (a applicationsService) HandlerGetApplicationById(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := a.applicationsRepository.FindById(id)

	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
