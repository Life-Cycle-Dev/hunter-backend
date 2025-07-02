package entity

import (
	"gorm.io/gorm"
	"time"
)

type OneTimePasswordType string

const (
	OneTimePasswordVerifyEmail OneTimePasswordType = "verify_email"
)

type OneTimePassword struct {
	ID        string              `json:"id" gorm:"type:varchar(255);primarykey"`
	UserId    string              `json:"user_id" gorm:"type:varchar(255);index"`
	Type      OneTimePasswordType `json:"type" gorm:"type:varchar(255)"`
	Code      string              `json:"code" gorm:"type:varchar(255)"`
	Ref       string              `json:"ref" gorm:"type:varchar(255)"`
	Revoke    bool                `json:"revoke" gorm:"default:false"`
	ExpiredAt time.Time           `json:"expired_at"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt gorm.DeletedAt      `json:"deleted_at" gorm:"index"`
}
