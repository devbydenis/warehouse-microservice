package request

type CreateWarehouseProductRequest struct {
	WarehouseID uint `json:"warehouse_id" validate:"required"`
	ProductID   uint `json:"product_id" validate:"required"`
	Stock       int  `json:"stock" validate:"required"`
}


