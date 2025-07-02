package authService

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hunter-backend/entity"
	"hunter-backend/util"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a authService) HandlerLogin(c *fiber.Ctx) error {
	var request LoginRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	user, err := a.userRepository.FindByEmail(a.encryptorRepository.Encrypt(request.Email))
	if err != nil {
		panic(err)
	}

	err = a.userRepository.CheckPassword(user.HashedPassword, request.Password)
	if err != nil {
		panic("invalid password")
	}

	if !user.IsEmailVerified {
		oneTimePasswordEnt, err := a.oneTimePasswordRepository.CreateOneTimePassword(user, entity.OneTimePasswordVerifyEmail)
		if err != nil {
			return err
		}

		err = a.notificationRepository.SendNotification(&entity.Notification{
			Type:    entity.NotificationEmail,
			Email:   user.Email,
			Title:   "OTP for E-mail Address verification on Hunter App",
			Content: fmt.Sprintf(util.GetEmailContent("verifyEmail"), a.encryptorRepository.Decrypt(user.Name), oneTimePasswordEnt.Code, oneTimePasswordEnt.Ref),
		})

		return c.JSON(fiber.Map{
			"uid":       user.ID,
			"verifyRef": oneTimePasswordEnt.Ref,
		})
	}

	jwtToken, err := a.jsonWebTokenRepository.GenerateToken(user)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(jwtToken)
}
