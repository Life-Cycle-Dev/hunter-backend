package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
	"hunter-backend/util"
)

type CreatePermissionRequest struct {
	Title   string `json:"title" validate:"required"`
	Mapping string `json:"mapping" validate:"required"`
}

func (p permissionService) HandlerCreatePermission(c *fiber.Ctx) error {
	var request CreatePermissionRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	exitingPermission, _ := p.permissionRepository.GetByMapping(request.Mapping)
	if exitingPermission.ID != "" {
		panic("permission with mapping already exists")
	}

	createdPermission, err := p.permissionRepository.CreatePermission(&entity.Permission{
		Title:   request.Title,
		Mapping: request.Mapping,
	})
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(createdPermission)
}
