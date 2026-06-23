package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}