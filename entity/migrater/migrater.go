package migrater

import (
	"gorm.io/gorm"
	"hunter-backend/entity"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Encryptor{},
		&entity.JsonWebToken{},
		&entity.OneTimePassword{},
		&entity.Notification{},
		&entity.Users{})
	if err != nil {
		return err
	}
	return nil
}
