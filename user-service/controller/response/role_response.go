package response

type RoleResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	CountUsers int64  `json:"count_users"`
}
