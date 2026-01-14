package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-warehouse/product-service/configs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type WarehouseClient struct {
	urlWarehouseService string
	httpClient          *http.Client
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

func NewWarehouseClient(cfg configs.Config) *WarehouseClient {
	return &WarehouseClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlWarehouseService: cfg.App.UrlWarehouseService,
	}
}

func (wc *WarehouseClient) GetProductStockAcrossWarehouses(ctx context.Context, productID uint) (int, error) {
	url := fmt.Sprintf("%s/api/v1/warehouse-products/detail/products/%d/total-stock", wc.urlWarehouseService, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 1: %v", err)
		return 0, err
	}

	resp, err := wc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 2: %v", err)
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 3: %v", err)
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 4: %s", string(body))
		return 0, errors.New("failed to get product stock across warehouses")
	}

	var stockResp WarehouseProductStockServiceResponse
	if err := json.Unmarshal(body, &stockResp); err != nil {
		log.Errorf("[WarehouseClient] GetProductStockAcrossWarehouses - 5: %v", err)
		return 0, err
	}

	return stockResp.Data.TotalStock, nil
}

func (wc *WarehouseClient) DeleteAllProductWarehouseProducts(ctx context.Context, productID uint) error {
	url := fmt.Sprintf("%s/api/v1/warehouse-products/detail/products/%d", wc.urlWarehouseService, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 1: %v", err)
		return err
	}

	resp, err := wc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 2: %v", err)
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 2: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] DeleteAllProductWarehouseProducts - 1: %s", string(body))
		return errors.New("failed to delete all product warehouse products")
	}

	return nil
}