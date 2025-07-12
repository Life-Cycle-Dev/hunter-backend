package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

type UpdateRoleRequest struct {
	Title         string   `json:"title" validate:"required"`
	Mapping       string   `json:"mapping" validate:"required"`
	PermissionIds []string `json:"permission_ids"`
}

func (p permissionService) HandlerUpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var request UpdateRoleRequest
	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	roleEnt, err := p.roleRepository.FindById(id)
	if err != nil {
		panic(err)
	}

	roleEnt.Title = request.Title
	roleEnt.Mapping = request.Mapping

	exitingRole, _ := p.roleRepository.FindByMapping(request.Mapping)
	if exitingRole.ID != "" && exitingRole.ID != id {
		panic("role with mapping already exists")
	}

	updatedEnt, err := p.roleRepository.UpdateRole(roleEnt)
	if err != nil {
		panic(err)
	}

	err = p.permissionRepository.UnlinkPermissionByRoleId(roleEnt.ID)

	permissionsEnt, err := p.permissionRepository.GetByIds(request.PermissionIds)
	if err != nil {
		panic(err)
	}

	var existingPermissionsIds []string
	for _, permission := range permissionsEnt {
		existingPermissionsIds = append(existingPermissionsIds, permission.ID)
	}

	_, err = p.permissionRepository.CreateRoleToPermission(roleEnt.ID, existingPermissionsIds)
	if err != nil {
		panic(err)
	}

	permissions, err := p.permissionRepository.GetByIds(request.PermissionIds)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"permissions": permissions,
		"role":        updatedEnt,
	})
}
