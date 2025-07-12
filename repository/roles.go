package repository

import (
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
	"strings"
)

type RoleRepository interface {
	CreateRole(ent *entity.Role, permissionIds []string) (*entity.Role, error)
	ListRoles(offset int, limit int, query string) ([]*entity.Role, int64, error)
	UpdateRole(ent *entity.Role) (*entity.Role, error)

	FindById(id string) (*entity.Role, error)
	FindByMapping(mapping string) (*entity.Role, error)
}

type roleRepository struct {
	db                   *gorm.DB
	config               config.AppConfig
	encryptorRepository  EncryptorRepository
	permissionRepository PermissionRepository
}

func (r roleRepository) UpdateRole(ent *entity.Role) (*entity.Role, error) {
	result := r.db.Updates(ent)
	if result.Error != nil {
		return nil, result.Error
	}
	return ent, nil
}

func (r roleRepository) FindByMapping(mapping string) (*entity.Role, error) {
	var ent entity.Role
	result := r.db.First(&ent, "mapping = ?", mapping)
	if result.Error != nil {
		return &entity.Role{}, result.Error
	}
	return &ent, nil
}

func (r roleRepository) ListRoles(offset int, limit int, query string) ([]*entity.Role, int64, error) {
	var roles []*entity.Role
	var total int64

	db := r.db.Model(&entity.Role{})

	if query != "" {
		likeQuery := "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(title) LIKE ? OR LOWER(mapping) LIKE ?", likeQuery, likeQuery)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r roleRepository) FindById(id string) (*entity.Role, error) {
	var ent entity.Role
	result := r.db.First(&ent, "id = ?", id)
	if result.Error != nil {
		return &entity.Role{}, result.Error
	}
	return &ent, nil
}

func (r roleRepository) CreateRole(ent *entity.Role, permissionIds []string) (*entity.Role, error) {
	id, err := r.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	ent.ID = id

	result := r.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	_, err = r.permissionRepository.CreateRoleToPermission(id, permissionIds)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

func ProvideRoleRepository(db *gorm.DB, config config.AppConfig) RoleRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	permissionRepository := ProvidePermissionRepository(db, config)
	return &roleRepository{
		db:                   db,
		config:               config,
		encryptorRepository:  encryptorRepository,
		permissionRepository: permissionRepository,
	}
}
