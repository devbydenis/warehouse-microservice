package model

import "time"

type Warehouse struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	Address   string     `json:"address" gorm:"type:text"`
	Photo     string     `json:"photo" gorm:"type:text"`
	Phone     string     `json:"phone" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	WarehouseProducts []WarehouseProduct `json:"warehouse_products" gorm:"foreignKey:WarehouseID"`
}
