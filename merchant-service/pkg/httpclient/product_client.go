package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-warehouse/merchant-service/configs"
	"micro-warehouse/merchant-service/pkg/jwt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type ProductClientInterface interface {
	GetProductByID(ctx context.Context, productID uint) (*ProductResponse, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*ProductResponse, error)
	GetProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]ProductResponse, error)
	HealthCheck(ctx context.Context) error
}

type ProductClient struct {
	urlProductService string
	httpClient        *http.Client
	config              configs.Config
}

func (p *ProductClient) generateInternalToken() (string, error) {
	return jwt.GenerateInternalToken(p.config)
}

// GetProductByBarcode implements ProductClientInterface.
func (p *ProductClient) GetProductByBarcode(ctx context.Context, barcode string) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products/barcode/%s", p.urlProductService, barcode)

	token, err := p.generateInternalToken()
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByBarcode - 1: %v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByBarcode - 2: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByBarcode - 3: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByBarcode - 4: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[ProductClient] - GetProductByBarcode - 5: %s", string(body))
		return nil, errors.New("failed to get product by barcode")
	}

	var productResponse ProductServiceResponse
	if err := json.Unmarshal(body, &productResponse); err != nil {
		log.Errorf("[ProductClient] - GetProductByBarcode - 6: %v", err)
		return nil, err
	}

	return &productResponse.Data, nil
}

// GetProductByID implements ProductClientInterface.
func (p *ProductClient) GetProductByID(ctx context.Context, productID uint) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products/%d", p.urlProductService, productID)

	token, err := p.generateInternalToken()
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByID - 1: %v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByID - 2: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByID - 3: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[ProductClient] - GetProductByID - 4: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[ProductClient] - GetProductByID - 5: %s", string(body))
		return nil, errors.New("failed to get product by id")
	}

	var productResponse ProductServiceResponse
	if err := json.Unmarshal(body, &productResponse); err != nil {
		log.Errorf("[ProductClient] - GetProductByID - 6: %v", err)
		return nil, err
	}

	return &productResponse.Data, nil
}

// GetProducts implements ProductClientInterface.
func (p *ProductClient) GetProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]ProductResponse, error) {
	url := fmt.Sprintf("%s/api/v1/products?page=%d&limit=%d&search=%s&sort_by=%s&sort_order=%s", p.urlProductService, page, limit, search, sortBy, sortOrder)

	token, err := p.generateInternalToken()
	if err != nil {
		log.Errorf("[ProductClient] - GetProducts - 1: %v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] - GetProducts - 2: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] - GetProducts - 3: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[ProductClient] - GetProducts - 4: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[ProductClient] - GetProducts - 5: %s", string(body))
		return nil, errors.New("failed to get products")
	}

	var productListResponse ProductListServiceResponse
	if err := json.Unmarshal(body, &productListResponse); err != nil {
		log.Errorf("[ProductClient] - GetProducts - 6: %v", err)
		return nil, err
	}

	return productListResponse.Data, nil
}

// HealthCheck implements ProductClientInterface.
func (p *ProductClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", p.urlProductService)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[ProductClient] - HealthCheck - 1: %v", err)
		return err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		log.Errorf("[ProductClient] - HealthCheck - 2: %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("product service is unhealthy")
	}

	return nil
}

type ProductResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	About     string `json:"about"`
	Price     int64  `json:"price"`
	Barcode   string `json:"barcode"`
	Thumbnail string `json:"thumbnail"`
	Category  struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
	} `json:"category"`
}

type ProductServiceResponse struct {
	Message string          `json:"message"`
	Data    ProductResponse `json:"data"`
	Error   string          `json:"error,omitempty"`
}

type ProductListServiceResponse struct {
	Message string            `json:"message"`
	Data    []ProductResponse `json:"data"`
	Error   string            `json:"error,omitempty"`
}

func NewProductClient(cfg configs.Config) ProductClientInterface {
	return &ProductClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlProductService: cfg.App.UrlProductService,
	}
}
