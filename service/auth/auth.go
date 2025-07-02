package authService

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/repository"
)

type AuthService interface {
	HandlerLogin(c *fiber.Ctx) error
	HandlerSignUp(c *fiber.Ctx) error
	HandlerVerifyEmail(c *fiber.Ctx) error
	HandlerGetUserInfo(c *fiber.Ctx) error
	HandlerRefreshAccessToken(c *fiber.Ctx) error
}

type authService struct {
	db                        *gorm.DB
	config                    config.AppConfig
	encryptorRepository       repository.EncryptorRepository
	userRepository            repository.UserRepository
	jsonWebTokenRepository    repository.JsonWebTokenRepository
	oneTimePasswordRepository repository.OneTimePasswordRepository
	notificationRepository    repository.NotificationRepository
}

func ProvideAuthService(db *gorm.DB, config config.AppConfig) AuthService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	userRepository := repository.ProvideUserRepository(db, config)
	jsonWebTokenRepository := repository.ProvideJsonWebTokenRepository(db, config)
	oneTimePasswordRepository := repository.ProvideOneTimePasswordRepository(db, config)
	notificationRepository := repository.ProvideNotificationRepository(db, config)
	return &authService{
		db:                        db,
		config:                    config,
		encryptorRepository:       encryptorRepository,
		userRepository:            userRepository,
		jsonWebTokenRepository:    jsonWebTokenRepository,
		oneTimePasswordRepository: oneTimePasswordRepository,
		notificationRepository:    notificationRepository,
	}
}
