package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"micro-warehouse/api-gateway/middleware"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userServiceURL string
	jwtConfig      middleware.JWTConfig
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role_name"`
}

type UserServiceResponse struct {
	User struct {
		UserID uint     `json:"user_id"`
		Email  string   `json:"email"`
		Role   []string `json:"role_name"`
	} `json:"user"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Role string `json:"role_name"`
	} `json:"user"`
}

func NewAuthController(userServiceURL string, jwtConfig middleware.JWTConfig) *AuthController {
	return &AuthController{
		userServiceURL: userServiceURL,
		jwtConfig:      jwtConfig,
	}
}

func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest LoginRequest
	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// validate request body
	if loginRequest.Email == "" || loginRequest.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Email and password are required",
		})
	}

	loginResp, err := ac.forwardLoginRequest(loginRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token, err := middleware.GenerateJWT(loginResp.UserID, loginResp.Email, loginResp.Role, ac.jwtConfig)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response := AuthResponse{
		Token: token,
		User: struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
			Role string `json:"role_name"`
		}{
			ID:    loginResp.UserID,
			Email: loginResp.Email,
			Role:  loginResp.Role,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    response,
	})
}

func (ac *AuthController) forwardLoginRequest(loginReq LoginRequest) (*LoginResponse, error) {
	log.Printf("STEP 1: Starting forwardLoginRequest")

	reqBody, err := json.Marshal(loginReq)
	log.Printf("STEP 2: Marshaled request: %s", string(reqBody))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ac.userServiceURL+"/api/v1/auth/login", bytes.NewBuffer(reqBody))
	log.Printf("STEP 3: Created request to %s", ac.userServiceURL+"/api/v1/auth/login")
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gateway", "warehouse-api-gateway")
	req.Header.Set("X-Internal-Request", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	log.Printf("STEP 4: Got response, err=%v", err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	log.Printf("STEP 5: Read body, status=%d, body=%s", resp.StatusCode, string(respBody))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Printf("STEP 6: Status != 200, returning error")
		return nil, err
	}

	var userServiceResp UserServiceResponse
	if err := json.Unmarshal(respBody, &userServiceResp); err != nil {
		return nil, err
	}

	rolesStr := ""
	if len(userServiceResp.User.Role) > 0 {
		rolesStr = userServiceResp.User.Role[0]
		for i := 1; i < len(userServiceResp.User.Role); i++ {
			rolesStr += "," + userServiceResp.User.Role[i]
		}
	}
	fmt.Println("USER SERVICE RESPONSE -> ", userServiceResp)
	
	loginResp := LoginResponse{
		UserID: userServiceResp.User.UserID,
		Email:  userServiceResp.User.Email,
		Role:   rolesStr,
	}

	return &loginResp, nil
}
