package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-warehouse/merchant-service/configs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type UserClientInterface interface {
	GetUserByID(ctx context.Context, userID uint) (*UserResponse, error)
}

type UserClient struct {
	urlUserService string
	httpClient     *http.Client
}

// GetUserByID implements UserClientInterface.
func (u *UserClient) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", u.urlUserService, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 1: %v", err)
		return nil, err
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 2: %v", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[UserClient] GetUserByID - 4: %s", string(body))
		return nil, errors.New("failed to get user by id")
	}

	var userResponse UserServiceResponse
	if err := json.Unmarshal(body, &userResponse); err != nil {
		log.Errorf("[UserClient] GetUserByID - 5: %v", err)
		return nil, err
	}

	return &userResponse.Data, nil
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Photo string `json:"photo"`
}

type UserServiceResponse struct {
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
	Error   string       `json:"error,omitempty"`
}

func NewUserCLient(cfg configs.Config) UserClientInterface {
	return &UserClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlUserService: cfg.App.UrlProductService,
	}
}
