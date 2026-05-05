package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-warehouse/product-service/configs"
	"micro-warehouse/product-service/pkg/jwt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type WarehouseClient struct {
	urlApiGateway string
	httpClient          *http.Client
	config              configs.Config
}

type WarehouseProductStockResponse struct {
	ProductID  uint `json:"product_id"`
	TotalStock int  `json:"total_stock"`
}

type WarehouseProductStockServiceResponse struct {
	Message string                        `json:"message"`
	Data    WarehouseProductStockResponse `json:"data"`
	Error   string                        `json:"error,omitempty"`
}

func (wc *WarehouseClient) generateInternalToken() (string, error) {
	return jwt.GenerateInternalToken(wc.config)
}

func NewWarehouseClient(cfg configs.Config) *WarehouseClient {
	return &WarehouseClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlApiGateway: cfg.App.UrlApiGateway,
		config:              cfg,
	}
}

func (wc *WarehouseClient) GetProductStockAcrossWarehouses(ctx context.Context, productID uint) (int, error) {
	url := fmt.Sprintf("%s/api/v1/warehouse-products/detail/products/%d/total-stock", wc.urlApiGateway, productID)

	token, err := wc.generateInternalToken()
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 1: %v", err)
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 2: %v", err)
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := wc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 3: %v", err)
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 4: %v", err)
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 5: %s", string(body))
		return 0, errors.New("failed to get product stock across warehouses")
	}

	var stockResp WarehouseProductStockServiceResponse
	if err := json.Unmarshal(body, &stockResp); err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 6: %v", err)
		return 0, err
	}

	return stockResp.Data.TotalStock, nil
}

func (wc *WarehouseClient) DeleteAllProductWarehouseProducts(ctx context.Context, productID uint) error {
	url := fmt.Sprintf("%s/api/v1/warehouse-products/detail/products/%d", wc.urlApiGateway, productID)

	token, err := wc.generateInternalToken()
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 1: %v", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 2: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := wc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 3: %v", err)
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 4: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 5: %s", string(body))
		return errors.New("failed to delete all product warehouse products")
	}

	return nil
}