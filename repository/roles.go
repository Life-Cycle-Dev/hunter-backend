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

	FindById(id string) (*entity.Role, error)
}

type roleRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
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

	roleToPermissions := make([]*entity.RoleToPermission, 0, len(permissionIds))
	for _, pid := range permissionIds {
		roleToPermissions = append(roleToPermissions, &entity.RoleToPermission{
			RoleId:       ent.ID,
			PermissionId: pid,
		})
	}

	if len(roleToPermissions) > 0 {
		if err := r.db.Create(&roleToPermissions).Error; err != nil {
			return nil, err
		}
	}

	return ent, nil
}

func ProvideRoleRepository(db *gorm.DB, config config.AppConfig) RoleRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &roleRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
