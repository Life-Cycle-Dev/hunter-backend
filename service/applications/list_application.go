package applicationsService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
	"strconv"
)

func (a applicationsService) HandlerListApplication(c *fiber.Ctx) error {
	permissions := c.Locals("permissions").([]string)
	err := util.CheckAllow(permissions, "application-view")
	if err != nil {
		panic(err)
	}

	pageStr := c.Query("page", "1")
	perPageStr := c.Query("perPage", "10")
	query := c.Query("query")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	apps, total, err := a.applicationsRepository.ListApplications(offset, perPage, query)
	if err != nil {
		panic(err)
	}

	totalPages := (int(total) + perPage - 1) / perPage

	return c.JSON(fiber.Map{
		"data": apps,
		"pagination": fiber.Map{
			"page":       page,
			"per_page":   perPage,
			"total":      total,
			"total_page": totalPages,
		},
		"query": query,
	})
}
