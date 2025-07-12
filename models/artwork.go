package models

import "time"

type Artwork struct {
	ID 				uint      `json:"id" gorm:"primaryKey"`
	Title 			string    `json:"title"`
	Caption 		string    `json:"caption"`
	ImageURL 		string    `json:"image_url"`
	CreatedAt 		time.Time `json:"created_at" gorm:"autoCreateTime"`
}