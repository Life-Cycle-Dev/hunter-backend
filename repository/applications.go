package repository

import (
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
)

type ApplicationsRepository interface {
	CreateApplication(applicationEnt *entity.Applications) (*entity.Applications, error)
}

type applicationsRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (a applicationsRepository) CreateApplication(applicationEnt *entity.Applications) (*entity.Applications, error) {
	id, err := a.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	applicationEnt.ID = id

	result := a.db.Create(applicationEnt)
	if result.Error != nil {
		return nil, result.Error
	}
	return applicationEnt, nil
}

func ProvideApplicationsRepository(db *gorm.DB, config config.AppConfig) ApplicationsRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &applicationsRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
