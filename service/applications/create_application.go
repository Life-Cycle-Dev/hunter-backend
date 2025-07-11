package applicationsService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
	"hunter-backend/util"
)

type CreateApplicationRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Active      bool   `json:"active"`
}

func (a applicationsService) HandlerCreateApplication(c *fiber.Ctx) error {
	var request CreateApplicationRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	ent, err := a.applicationsRepository.CreateApplication(&entity.Applications{
		Title:       request.Title,
		Description: request.Description,
		ImageUrl:    request.ImageUrl,
	})
	if err != nil {
		panic(err)
	}

	ent.Active = request.Active
	updatedEnt, err := a.applicationsRepository.UpdateApplication(ent)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusCreated).JSON(updatedEnt)
}
