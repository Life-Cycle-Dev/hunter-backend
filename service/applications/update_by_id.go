package applicationsService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

type UpdateApplicationRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Active      bool   `json:"active"`
}

func (a applicationsService) HandlerUpdateApplicationById(c *fiber.Ctx) error {
	id := c.Params("id")
	var request UpdateApplicationRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	applicationEnt, err := a.applicationsRepository.FindById(id)
	if err != nil {
		panic(err)
	}

	applicationEnt.Title = request.Title
	applicationEnt.Description = request.Description
	applicationEnt.ImageUrl = request.ImageUrl
	applicationEnt.Active = request.Active

	updatedEnt, err := a.applicationsRepository.UpdateApplication(applicationEnt)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedEnt)
}
