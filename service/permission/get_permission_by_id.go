package permissionService

import "github.com/gofiber/fiber/v2"

func (p permissionService) HandlerGetPermissionById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := p.permissionRepository.GetById(id)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
