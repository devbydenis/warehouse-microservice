package controller

import (
	request "micro-warehouse/user-service/controller/request"
	response "micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/pkg/conv"
	"micro-warehouse/user-service/pkg/validator"
	"micro-warehouse/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthControllerInterface interface {
	Login(c *fiber.Ctx) error
}

type AuthController struct {
	AuthService usecase.UserUsecaseInterface
}

// Login implements AuthControllerInterface.
// @Summary User Login
// @Description Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login Request Body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (a *AuthController) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		log.Errorf("[AuthController.Login] Login - 1: %v", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(loginRequest); err != nil {
		log.Errorf("[AuthController.Login] Login - 2: %v", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := a.AuthService.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		log.Errorf("[AuthController.Login] Login - 3: %v", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if user == nil {
		// log.Errorf("[AuthController.Login] Login - 4: %v", err.Error())
		log.Errorf("[AuthController.Login] Login - 4: %v", "user didn't exist")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	isValidPassword := conv.CheckPasswordHash(loginRequest.Password, user.Password)
	if !isValidPassword {
		// log.Errorf("[AuthController.Login] Login - 5: %v", err.Error())
		log.Errorf("[AuthController.Login] Login - 5: %v", "email or password not valid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	loginResp := response.LoginResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   roles, // disini kita buat seperti ini karna akan di mapping ulang lagi di api-gateway
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    loginResp,
	})
}

func NewAuthController(authService usecase.UserUsecaseInterface) AuthControllerInterface {
	return &AuthController{
		AuthService: authService,
	}
}
