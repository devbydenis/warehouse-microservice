package request

type CreateCategoryRequest struct {
	Name    string `json:"name" validate:"required"`
	Tagline string `json:"tagline" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}

type GetAllCategoryRequest struct {
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Search    string `query:"search"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
}

