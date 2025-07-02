package authService

import (
	"github.com/gofiber/fiber/v2"
	"hunter-backend/util"
	"time"
)

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
	Ref   string `json:"ref" validate:"required"`
}

func (a authService) HandlerVerifyEmail(c *fiber.Ctx) error {
	var request VerifyEmailRequest

	err := util.ValidateRequest(c, &request)
	if err != nil {
		panic(err)
	}

	user, err := a.userRepository.FindByEmail(a.encryptorRepository.Encrypt(request.Email))
	if err != nil {
		panic(err)
	}

	otpEnt, err := a.oneTimePasswordRepository.GetOneTimePassword(user, request.Ref)
	if err != nil || otpEnt == nil {
		panic("no one time password found")
	}

	if otpEnt.ExpiredAt.Before(time.Now()) {
		panic("otp is expired")
	}

	if otpEnt.Code != request.Code {
		panic("otp is invalid")
	}

	if otpEnt.Revoke {
		panic("otp is revoked")
	}

	otpEnt.Revoke = true
	_, err = a.oneTimePasswordRepository.UpdateOneTimePassword(otpEnt)
	if err != nil {
		panic(err)
	}

	user.IsEmailVerified = true
	updatedUser, err := a.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	jwtToken, err := a.jsonWebTokenRepository.GenerateToken(updatedUser)
	if err != nil {
		panic(err)
	}

	return c.Status(fiber.StatusOK).JSON(jwtToken)
}
