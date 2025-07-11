package repository

import (
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
	"strings"
)

type PermissionRepository interface {
	CreatePermission(ent *entity.Permission) (*entity.Permission, error)
	UpdatePermission(ent *entity.Permission) (*entity.Permission, error)
	ListPermissions(offset int, limit int, query string) ([]*entity.Permission, int64, error)

	GetByMapping(mapping string) (*entity.Permission, error)
	GetById(id string) (*entity.Permission, error)
}

type permissionRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (p permissionRepository) UpdatePermission(ent *entity.Permission) (*entity.Permission, error) {
	result := p.db.Updates(ent)
	if result.Error != nil {
		return nil, result.Error
	}
	return ent, nil
}

func (p permissionRepository) ListPermissions(offset int, limit int, query string) ([]*entity.Permission, int64, error) {
	var permissions []*entity.Permission
	var total int64

	db := p.db.Model(&entity.Permission{})

	if query != "" {
		likeQuery := "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(title) LIKE ? OR LOWER(mapping) LIKE ?", likeQuery, likeQuery)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (p permissionRepository) GetById(id string) (*entity.Permission, error) {
	var ent entity.Permission
	result := p.db.First(&ent, "id = ?", id)
	if result.Error != nil {
		return &entity.Permission{}, result.Error
	}
	return &ent, nil
}

func (p permissionRepository) GetByMapping(mapping string) (*entity.Permission, error) {
	var ent entity.Permission
	result := p.db.First(&ent, "mapping = ?", mapping)
	if result.Error != nil {
		return &entity.Permission{}, result.Error
	}
	return &ent, nil
}

func (p permissionRepository) CreatePermission(ent *entity.Permission) (*entity.Permission, error) {
	id, err := p.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	ent.ID = id

	result := p.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}
	return ent, nil
}

func ProvidePermissionRepository(db *gorm.DB, config config.AppConfig) PermissionRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &permissionRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
