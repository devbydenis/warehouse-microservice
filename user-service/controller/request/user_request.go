package request

type AssignUserToRoleRequest struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Photo    string `json:"photo" validate:"required"`
}

type GetAllUsersRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search string `query:"search" validate:"omitempty"`
	SortBy string `query:"sort_by" validate:"omitempty,oneof=id name email created_at"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Phone string `json:"phone"`
	Photo string `json:"photo"`
}