package request

type CreateWarehouseProductRequest struct {
	ProductID   uint `json:"product_id" validate:"required"`
	Stock       int  `json:"stock" validate:"required"`
}


