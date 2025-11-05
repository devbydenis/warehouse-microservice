package response

import "micro-warehouse/product-service/pkg/pagination"

type CategoryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Tagline      string `json:"tagline"`
	Photo        string `json:"photo"`
	CountProduct int    `json:"count_product"`
}

type GetAllCategoryResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}
