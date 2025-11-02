package request

type CreateProductRequest struct {
	Name       string `json:"name" validate:"required"`
	Barcode    string `json:"barcode" validate:"required"`
	Price      float64    `json:"price" validate:"required"`
	About      string `json:"about" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	Thumbnail  string `json:"thumbnail" validate:"required"`
	IsPopular  bool   `json:"is_popular" validate:"required"`
}

type GetAllProductRequest struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}


