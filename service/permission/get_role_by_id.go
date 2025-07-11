package permissionService

import "github.com/gofiber/fiber/v2"

func (p permissionService) HandlerGetRoleById(c *fiber.Ctx) error {
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
