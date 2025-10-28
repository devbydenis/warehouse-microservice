package response

import "micro-warehouse/user-service/pkg/pagination"

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	RoleName string `json:"role_name"`
}

type GetAllUsersResponse struct {
	Users      []UserResponse                `json:"users"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}

type UserRoleResponse struct {
	ID     uint         `json:"id"`
	UserID uint         `json:"user_id"`
	RoleID uint         `json:"role_id"`
	User   UserResponse `json:"user"`
	Role   RoleResponse `json:"role"`
}

type GetAllUserRolesResponse struct {
	UserRoles  []UserRoleResponse            `json:"user_roles"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}