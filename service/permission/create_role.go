package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
	"hunter-backend/util"
)

type CreateRoleRequest struct {
	Title         string   `json:"title" validate:"required"`
	Mapping       string   `json:"mapping" validate:"required"`
	PermissionIds []string `json:"permission_ids"`
}

func (p permissionService) HandlerCreateRole(c *fiber.Ctx) error {
	var request CreateRoleRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	permissionsEnt, err := p.permissionRepository.GetByIds(request.PermissionIds)
	if err != nil {
		panic(err)
	}

	var existingPermissionsIds []string
	for _, permission := range permissionsEnt {
		existingPermissionsIds = append(existingPermissionsIds, permission.ID)
	}

	roleEnt := &entity.Role{
		Title:   request.Title,
		Mapping: request.Mapping,
	}

	createdEnt, err := p.roleRepository.CreateRole(roleEnt, existingPermissionsIds)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"permissions": permissionsEnt,
		"role":        createdEnt,
	})
}
