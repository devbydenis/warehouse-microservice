package pagination

import "math"

type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	Limit        int `json:"limit"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

func CalculatePagination(page, limit, totalRecords int) PaginationResponse {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalRecords: int64(totalRecords),
		Limit:        limit,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}
