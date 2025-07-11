package repository

import (
	"gorm.io/gorm"
	"hunter-backend/di/config"
)

type RoleRepository interface {
}

type roleRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func ProvideRoleRepository(db *gorm.DB, config config.AppConfig) RoleRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &roleRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
