package domain

import (
	"time"

	"gorm.io/gorm"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type GetProfileResponse struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	IsPremium bool   `json:"is_premium" gorm:"default:false"`

	// Profile
	Profile *Profile `json:"profile" gorm:"foreignKey:UserID"`
}
