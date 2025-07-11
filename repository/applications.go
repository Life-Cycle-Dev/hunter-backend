package repository

import (
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
	"strings"
)

type ApplicationsRepository interface {
	CreateApplication(applicationEnt *entity.Applications) (*entity.Applications, error)

	ListApplications(offset int, limit int, query string) ([]*entity.Applications, int64, error)
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

func (a applicationsRepository) ListApplications(offset, limit int, query string) ([]*entity.Applications, int64, error) {
	var applications []*entity.Applications
	var total int64

	db := a.db.Model(&entity.Applications{})

	if query != "" {
		likeQuery := "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", likeQuery, likeQuery)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

func ProvideApplicationsRepository(db *gorm.DB, config config.AppConfig) ApplicationsRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &applicationsRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
