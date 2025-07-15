package userService

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
)

type UpdateUpdateUserRequest struct {
	Email           string `json:"email" validate:"required"`
	Name            string `json:"name" validate:"required"`
	RoleId          string `json:"role_id" validate:"required"`
	IsDeveloper     bool   `json:"is_developer"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

func (u userService) HandlerUpdateUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	var request UpdateUpdateUserRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	user, err := u.userRepository.FindById(id)
	if err != nil {
		panic(err)
	}

	existingRole, err := u.roleRepository.FindById(request.RoleId)
	if existingRole.ID == "" {
		panic(errors.New("role not found"))
	}

	user.Email = u.encryptorRepository.Encrypt(request.Email)
	user.Name = u.encryptorRepository.Encrypt(request.Name)
	user.RoleId = request.RoleId
	user.IsDeveloper = request.IsDeveloper
	user.IsEmailVerified = request.IsEmailVerified

	updatedUser, err := u.userRepository.UpdateUser(user)
	if err != nil {
		panic(err)
	}

	role, err := u.roleRepository.FindById(updatedUser.RoleId)
	if err != nil {
		panic(err)
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
