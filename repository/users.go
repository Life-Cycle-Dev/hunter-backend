package repository

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"hunter-backend/di/config"
	"hunter-backend/entity"
)

type UserRepository interface {
	CreateUser(ent *entity.Users) (*entity.Users, error)
	UpdateUser(ent *entity.Users) (*entity.Users, error)
	SignUpWithPassword(ent *entity.Users, password string) (*entity.Users, error)

	CheckPassword(hashedPassword entity.EncryptedField, plainPassword string) error
	FindByEmail(email entity.EncryptedField) (*entity.Users, error)
	FindById(id string) (*entity.Users, error)
}

type userRepository struct {
	db                  *gorm.DB
	config              config.AppConfig
	encryptorRepository EncryptorRepository
}

func (u userRepository) FindById(id string) (*entity.Users, error) {
	var ent entity.Users
	result := u.db.First(&ent, "id = ?", id)
	if result.Error != nil {
		return &entity.Users{}, result.Error
	}
	return &ent, nil
}

func (u userRepository) CheckPassword(hashedPassword entity.EncryptedField, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
}

func (u userRepository) FindByEmail(email entity.EncryptedField) (*entity.Users, error) {
	var ent entity.Users
	result := u.db.First(&ent, "email = ?", email)
	if result.Error != nil {
		return &entity.Users{}, result.Error
	}
	return &ent, nil
}

func (u userRepository) CreateUser(ent *entity.Users) (*entity.Users, error) {
	id, err := u.encryptorRepository.GeneratePassphrase(20)
	if err != nil {
		return nil, err
	}
	ent.ID = id

	result := u.db.Create(ent)
	if result.Error != nil {
		return nil, result.Error
	}

	return ent, nil
}

func (u userRepository) UpdateUser(ent *entity.Users) (*entity.Users, error) {
	result := u.db.Updates(ent)
	if result.Error != nil {
		return nil, result.Error
	}
	return ent, nil
}

func (u userRepository) SignUpWithPassword(ent *entity.Users, password string) (*entity.Users, error) {
	_, err := u.FindByEmail(ent.Email)

	if err == nil {
		return nil, errors.New("email already in use")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	ent.HashedPassword = hashedPassword

	createdUser, err := u.CreateUser(ent)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func ProvideUserRepository(db *gorm.DB, config config.AppConfig) UserRepository {
	encryptorRepository := ProvideEncryptorRepository(db, config)
	return &userRepository{
		db:                  db,
		config:              config,
		encryptorRepository: encryptorRepository,
	}
}
