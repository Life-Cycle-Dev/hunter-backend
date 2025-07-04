package entity

import (
	"gorm.io/gorm"
	"time"
)

type Applications struct {
	ID          string         `json:"id" gorm:"type:varchar(255);primary_key"`
	Title       string         `json:"title" gorm:"type:varchar(255);" validate:"required"`
	Description string         `json:"description" gorm:"type:varchar(255);"`
	ImageUrl    string         `json:"image_url" gorm:"type:varchar(255);"`
	Active      bool           `json:"active" gorm:"type:boolean;default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
