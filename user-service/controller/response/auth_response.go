package response

type LoginResponse struct {
	UserID uint     `json:"user_id"`
	Email  string   `json:"email"`
	Role   []string `json:"role_name"`
}
