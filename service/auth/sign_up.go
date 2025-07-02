package authService

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
	"hunter-backend/util"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a authService) HandlerSignUp(c *fiber.Ctx) error {
	var request SignUpRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	createdUser, err := a.userRepository.SignUpWithPassword(&entity.Users{
		Name:  a.encryptorRepository.Encrypt(request.Name),
		Email: a.encryptorRepository.Encrypt(request.Email),
	}, request.Password)
	if err != nil {
		panic(err)
	}

	oneTimePasswordEnt, err := a.oneTimePasswordRepository.CreateOneTimePassword(createdUser, entity.OneTimePasswordVerifyEmail)
	if err != nil {
		return err
	}

	err = a.notificationRepository.SendNotification(&entity.Notification{
		Type:    entity.NotificationEmail,
		Email:   createdUser.Email,
		Title:   "OTP for E-mail Address verification on Hunter App",
		Content: fmt.Sprintf(util.GetEmailContent("verifyEmail"), a.encryptorRepository.Decrypt(createdUser.Name), oneTimePasswordEnt.Code, oneTimePasswordEnt.Ref),
	})

	return c.JSON(fiber.Map{
		"uid":       createdUser.ID,
		"verifyRef": oneTimePasswordEnt.Ref,
	})
}
