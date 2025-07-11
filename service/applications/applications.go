package applicationsService

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/repository"
)

type ApplicationsService interface {
	HandlerCreateApplication(c *fiber.Ctx) error
	HandlerListApplication(c *fiber.Ctx) error
	HandlerGetApplicationById(c *fiber.Ctx) error
	HandlerUpdateApplicationById(c *fiber.Ctx) error
}

type applicationsService struct {
	db                     *gorm.DB
	config                 config.AppConfig
	encryptorRepository    repository.EncryptorRepository
	applicationsRepository repository.ApplicationsRepository
}

func ProvideApplicationsService(db *gorm.DB, config config.AppConfig) ApplicationsService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	applicationsRepository := repository.ProvideApplicationsRepository(db, config)
	return &applicationsService{
		db:                     db,
		config:                 config,
		encryptorRepository:    encryptorRepository,
		applicationsRepository: applicationsRepository,
	}
}
