package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

func (p permissionService) HandlerGetPermissionById(c *fiber.Ctx) error {
	permissions := c.Locals("permissions").([]string)
	err := util.CheckAllow(permissions, "permission-view")
	if err != nil {
		panic(err)
	}

	id := c.Params("id")

	result, err := p.permissionRepository.GetById(id)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
