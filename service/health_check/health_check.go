package healthCheckService

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/repository"
)

type HealthCheckService interface {
	HandlerGetRouter(c *fiber.Ctx) error
}

type healthCheckService struct {
	encryptorRepository repository.EncryptorRepository
	db                  *gorm.DB
	config              config.AppConfig
}

func (h healthCheckService) HandlerGetRouter(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"healthy": true,
	})
}

func ProvideHealthCheckService(db *gorm.DB, config config.AppConfig) HealthCheckService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	return &healthCheckService{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
