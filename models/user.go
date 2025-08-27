package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `gorm:"unique" json:"username"`
	Password     string    `json: "password"`
	RefreshToken string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
