package request

type CreateMerchantProductRequest struct {
	ProductID   uint `json:"product_id" validate:"required"`
	WarehouseID uint `json:"warehouse_id" validate:"required"`
	MerchantID  uint `json:"merchant_id" validate:"required"`
	Stock       int  `json:"stock" validate:"required"`
}

type GetMerchantProductRequest struct {
	Page       int    `query:"page" validate:"omitempty,min=1"`
	Limit      int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search     string `query:"search" validate:"omitempty"`
	SortBy     string `query:"sort_by" validate:"omitempty,oneof=id product_id warehouse_id stock created_at"`
	SortOrder  string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	MerchantID uint   `query:"merchant_id" validate:"omitempty"`
	ProductID  uint   `query:"product_id" validate:"omitempty"`
	KeeperID   uint   `query:"keeper_id" validate:"omitempty"`
}
