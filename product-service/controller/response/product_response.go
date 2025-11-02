package response

import "micro-warehouse/product-service/pkg/pagination"

type ProductResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Barcode    string           `json:"barcode"`
	Price      int              `json:"price"`
	About      string           `json:"about"`
	CategoryID uint             `json:"category_id"`
	Thumbnail  string           `json:"thumbnail"`
	IsPopular  bool             `json:"is_popular"`
	Category   CategoryResponse `json:"category"`
}

type GetAllProductResponse struct {
	Products   []ProductResponse             `json:"products"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}
