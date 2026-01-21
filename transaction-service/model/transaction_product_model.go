package model

import (
	"time"

	"gorm.io/gorm"
)

type TransactionProduct struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID     uint           `json:"product_id" gorm:"type:bigint;not null"`
	Quantity      int64          `json:"quantity" gorm:"type:bigint;not null"`
	Price         int64          `json:"price" gorm:"type:bigint;not null"`
	Subtotal      int64          `json:"sub_total" gorm:"type:bigint;not null"`
	TransactionID uint           `json:"transaction_id" gorm:"type:bigint;not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Virtual field for Response
	ProductName          string `json:"product_name" gorm:"-"`
	ProductPhoto         string `json:"product_photo" gorm:"-"`
	ProductAbout         string `json:"product_about" gorm:"-"`
	ProductCategoryID    uint `json:"product_category_id" gorm:"-"`
	ProductCategoryName  string `json:"product_category_name" gorm:"-"`
	ProductCategoryPhoto string `json:"product_category_photo" gorm:"-"`

	// Relationships
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:ID"`
}
