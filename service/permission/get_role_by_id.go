package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

func (p permissionService) HandlerGetRoleById(c *fiber.Ctx) error {
	userPermissions := c.Locals("permissions").([]string)
	err := util.CheckAllow(userPermissions, "role-view")
	if err != nil {
		panic(err)
	}

	id := c.Params("id")

	role, err := p.roleRepository.FindById(id)
	if err != nil {
		panic(err)
	}

	permissions, err := p.permissionRepository.GetByRoleId(role.ID)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"role":        role,
		"permissions": permissions,
	})
}
