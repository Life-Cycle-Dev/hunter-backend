package entity

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	ID              string         `json:"id" gorm:"type:varchar(255);primarykey"`
	Name            EncryptedField `json:"name" gorm:"type:varbinary(512)" validate:"required"`
	Email           EncryptedField `json:"email" gorm:"type:varbinary(512)" validate:"required"`
	HashedPassword  EncryptedField `json:"hashed_password" gorm:"type:varbinary(512)" validate:"required"`
	IsEmailVerified bool           `json:"is_email_verified" gorm:"type:boolean;default:false"`
	IsDeveloper     bool           `json:"is_developer" gorm:"type:boolean;default:false"`
	PermissionId    string         `json:"permission_id" gorm:"type:varchar(255);"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	IsDeveloper     bool   `json:"is_developer"`
	IsEmailVerified bool   `json:"is_email_verified"`
	CreatedAt       string `json:"created_at"`
}

func (e *Users) ToResponse(decrypt func(EncryptedField) string) UserResponse {
	return UserResponse{
		ID:              e.ID,
		Name:            decrypt(e.Name),
		Email:           decrypt(e.Email),
		IsDeveloper:     e.IsDeveloper,
		IsEmailVerified: e.IsEmailVerified,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
	}
}
