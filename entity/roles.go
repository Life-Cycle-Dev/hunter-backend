package entity

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID        string         `json:"id" gorm:"type:varchar(255);primary_key"`
	Title     string         `json:"title" gorm:"type:varchar(255);"`
	Mapping   string         `json:"mapping" gorm:"type:varchar(255);"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
