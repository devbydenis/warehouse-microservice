package model

import "time"

type Category struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	Tagline   string     `json:"tagline" gorm:"type:varchar(100);uniqueIndex"`
	Photo     string     `json:"photo" gorm:"type:text"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Products []Product `json:"products" gorm:"foreignKey:CategoryID"`
}
