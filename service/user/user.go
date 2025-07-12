package userService

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/repository"
)

type UserService interface {
	HandlerListUser(c *fiber.Ctx) error
}

type userService struct {
	db                   *gorm.DB
	config               config.AppConfig
	encryptorRepository  repository.EncryptorRepository
	permissionRepository repository.PermissionRepository
	roleRepository       repository.RoleRepository
	userRepository       repository.UserRepository
}

func ProvideUserService(db *gorm.DB, config config.AppConfig) UserService {
	encryptorRepository := repository.ProvideEncryptorRepository(db, config)
	permissionRepository := repository.ProvidePermissionRepository(db, config)
	roleRepository := repository.ProvideRoleRepository(db, config)
	userRepository := repository.ProvideUserRepository(db, config)
	return &userService{
		db:                   db,
		config:               config,
		encryptorRepository:  encryptorRepository,
		permissionRepository: permissionRepository,
		roleRepository:       roleRepository,
		userRepository:       userRepository,
	}
}
