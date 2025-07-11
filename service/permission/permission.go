package permissionService

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/repository"
)

type PermissionService interface {
	HandlerCreatePermission(c *fiber.Ctx) error
	HandlerGetPermissionById(c *fiber.Ctx) error
	HandlerListPermission(c *fiber.Ctx) error
	HandlerUpdatePermission(c *fiber.Ctx) error
}

type permissionService struct {
	db                   *gorm.DB
	config               config.AppConfig
	encryptorRepository  repository.EncryptorRepository
	permissionRepository repository.PermissionRepository
}

func ProvidePermissionService(db *gorm.DB, config config.AppConfig) PermissionService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	permissionRepository := repository.ProvidePermissionRepository(db, config)
	return &permissionService{
		db:                   db,
		config:               config,
		encryptorRepository:  encryptorRepository,
		permissionRepository: permissionRepository,
	}
}
