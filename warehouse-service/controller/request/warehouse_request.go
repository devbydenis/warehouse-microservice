package request

type CreateWarehouseRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}

type GetAllWarehouseRequest struct {
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search    string `query:"search" validate:"omitempty"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=id name address phone created_at"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

