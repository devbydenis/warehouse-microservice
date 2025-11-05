package model

import "time"

type WarehouseProduct struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	WarehouseID uint       `json:"warehouse_id" gorm:"not null;index"`
	ProductID   uint       `json:"product_id" gorm:"not null;index"`
	Stock       int        `json:"stock" gorm:"not null;default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`

	Warehouse 	Warehouse 	`json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}
