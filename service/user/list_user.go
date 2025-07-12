package userService

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (u userService) HandlerListUser(c *fiber.Ctx) error {
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

	users, total, err := u.userRepository.ListUser(offset, perPage, query)
	if err != nil {
		panic(err)
	}

	roleCache := make(map[string]string)
	var usersResponse []fiber.Map

	for _, user := range users {
		roleName := "User"

		if cachedTitle, ok := roleCache[user.RoleId]; ok {
			roleName = cachedTitle
		} else {
			role, err := u.roleRepository.FindById(user.RoleId)
			if err == nil {
				roleName = role.Title
			}
			roleCache[user.RoleId] = roleName
		}

		usersResponse = append(usersResponse, fiber.Map{
			"id":                user.ID,
			"name":              u.encryptorRepository.Decrypt(user.Name),
			"email":             u.encryptorRepository.Decrypt(user.Email),
			"is_email_verified": user.IsEmailVerified,
			"is_developer":      user.IsDeveloper,
			"role":              roleName,
		})
	}

	totalPages := (int(total) + perPage - 1) / perPage

	return c.JSON(fiber.Map{
		"data": usersResponse,
		"pagination": fiber.Map{
			"page":       page,
			"per_page":   perPage,
			"total":      total,
			"total_page": totalPages,
		},
		"query": query,
	})
}
