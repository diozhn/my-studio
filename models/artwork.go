package models

import "time"

type Artwork struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	ImageURL  string    `json:"image_url"`
	Likes     int       `json:"likes"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
