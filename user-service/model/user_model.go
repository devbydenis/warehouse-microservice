package model

import "time"

type User struct {
	ID       uint    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    string `json:"photo"`
	Phone    string `json:"phone"`
	Roles    []Role `gorm:"many2many:user_role;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
