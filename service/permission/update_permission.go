package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

type UpdatePermissionRequest struct {
	Title   string `json:"title" validate:"required"`
	Mapping string `json:"mapping" validate:"required"`
}

func (p permissionService) HandlerUpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")

	var request UpdatePermissionRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	permissionEnt, err := p.permissionRepository.GetById(id)
	if err != nil {
		panic(err)
	}

	permissionEnt.Title = request.Title
	permissionEnt.Mapping = request.Mapping

	exitingPermission, _ := p.permissionRepository.GetByMapping(request.Mapping)
	if exitingPermission.ID != "" && exitingPermission.ID != id {
		panic("permission with mapping already exists")
	}

	updatedEnt, err := p.permissionRepository.UpdatePermission(permissionEnt)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedEnt)
}
