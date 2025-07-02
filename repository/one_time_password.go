package repository

import (
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
	"time"
)

type OneTimePasswordRepository interface {
	CreateOneTimePassword(user *entity.Users, otpType entity.OneTimePasswordType) (*entity.OneTimePassword, error)
	GetOneTimePassword(user *entity.Users, ref string) (*entity.OneTimePassword, error)
	UpdateOneTimePassword(ent *entity.OneTimePassword) (*entity.OneTimePassword, error)
}

type oneTimePasswordRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (o oneTimePasswordRepository) UpdateOneTimePassword(ent *entity.OneTimePassword) (*entity.OneTimePassword, error) {
	result := o.db.Model(ent).Updates(ent)
	if result.Error != nil {
		return nil, result.Error
	}
	return ent, nil
}

func (o oneTimePasswordRepository) GetOneTimePassword(user *entity.Users, ref string) (*entity.OneTimePassword, error) {
	var ent entity.OneTimePassword
	result := o.db.First(&ent, "user_id = ? AND ref = ?", user.ID, ref)
	if result.Error != nil {
		return &entity.OneTimePassword{}, result.Error
	}
	return &ent, nil
}

func (o oneTimePasswordRepository) CreateOneTimePassword(user *entity.Users, otpType entity.OneTimePasswordType) (*entity.OneTimePassword, error) {
	id, err := o.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}

	code, err := o.encryptorRepository.GenerateNumericCode(6)
	if err != nil {
		return nil, err
	}

	ref, err := o.encryptorRepository.GenerateReferenceCode(4)
	if err != nil {
		return nil, err
	}

	ent := &entity.OneTimePassword{
		ID:        id,
		UserId:    user.ID,
		Type:      otpType,
		Code:      code,
		Ref:       ref,
		ExpiredAt: time.Now().Add(15 * time.Minute),
	}

	result := o.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	return ent, nil
}

func ProvideOneTimePasswordRepository(db *gorm.DB, config config.AppConfig) OneTimePasswordRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &oneTimePasswordRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
