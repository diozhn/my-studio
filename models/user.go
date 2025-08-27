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
	Email        string    `json:"email" gorm:"unique"`
	SuperUser    bool      `json:"superuser" gorm:"default:false"`
	GoogleID     string    `gorm:"uniqueIndex"`
	InstagramID  string    `gorm:"uniqueIndex"`
	TwitterID    string    `gorm:"uniqueIndex"`
	RefreshToken string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
