package userService

import "github.com/gofiber/fiber/v2"

func (u userService) HandlerGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := u.userRepository.FindById(id)
	if err != nil {
		panic(err)
	}

	role, err := u.roleRepository.FindById(user.RoleId)
	if err != nil {
		role, err = u.roleRepository.FindByMapping("user")
		if err != nil {
			panic(err)
		}
	}

	permissions, err := u.permissionRepository.GetByRoleId(role.ID)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":                user.ID,
		"name":              u.encryptorRepository.Decrypt(user.Name),
		"email":             u.encryptorRepository.Decrypt(user.Email),
		"role":              role,
		"role_id":           role.ID,
		"permissions":       permissions,
		"is_email_verified": user.IsEmailVerified,
		"is_developer":      user.IsDeveloper,
		"created_at":        user.CreatedAt,
		"updated_at":        user.UpdatedAt,
	})
}
