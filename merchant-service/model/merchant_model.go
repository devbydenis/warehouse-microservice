package model

import "time"

type Merchant struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"type:varchar(100);not null"`
	Address   string     `json:"address" gorm:"type:text"`
	Photo     string     `json:"photo"`
	Phone     string     `json:"phone"`
	KeeperID  uint       `json:"keeper_id" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	MerchantProducts []MerchantProduct `json:"merchant_products" gorm:"foreignKey:MerchantID"`
}
