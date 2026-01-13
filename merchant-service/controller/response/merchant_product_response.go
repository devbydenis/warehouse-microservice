package response

import (
	"micro-warehouse/merchant-service/pkg/pagination"
	"time"
)

type MerchantProductResponse struct {
	ID          uint      `json:"id"`
	MerchantID  uint      `json:"merchant_id"`
	ProductID   uint      `json:"product_id"`
	WarehouseID uint      `json:"warehouse_id"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Merchant MerchantResponse `json:"merchant,omitempty"`
}

type MerchantProduct struct {
	ID                   uint   `json:"id"`
	MerchantID           uint   `json:"merchant_id"`
	ProductID            uint   `json:"product_id"`
	ProductName          string `json:"product_name"`
	ProductAbout         string `json:"product_about"`
	ProductPhoto         string `json:"product_photo"`
	ProductPrice         int    `json:"product_price"`
	ProductCategory      string `json:"product_category"`
	ProductCategoryPhoto string `json:"product_category_photo"`
	Stock                int    `json:"stock"`
	WarehouseID          uint   `json:"warehouse_id"`
	WarehouseName        string `json:"warehouse_name"`
	WarehousePhoto       string `json:"warehouse_photo"`
	WarehousePhone       string `json:"warehouse_phone"`
}

type ProductTotalResponse struct {
	ProductID  uint `json:"product_id"`
	TotalStock int  `json:"total_stock"`
}

type GetAllMerchantProductsResponse struct {
	MerchantProducts []MerchantProduct             `json:"merchant_products"`
	Pagination       pagination.PaginationResponse `json:"pagination"`
}
