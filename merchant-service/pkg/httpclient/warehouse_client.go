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

type WarehouseClientInterface interface {
	GetWarehouseByID(ctc context.Context, warehouseID uint) (*WarehouseResponse, error)
	GetWarehouseProductStock(ctc context.Context, warehouseID, productID uint) (*WarehouseProductStockResponse, error)
}

type WarehouseCLient struct {
	urlWarehouseService string
	httpClient          *http.Client
}

// GetWarehouseByID implements WarehouseClientInterface.
func (w *WarehouseCLient) GetWarehouseByID(ctx context.Context, warehouseID uint) (*WarehouseResponse, error) {
	url := fmt.Sprintf("%s/api/v1/warehouses/%d", w.urlWarehouseService, warehouseID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseByID - 1: %v", err)
		return nil, err
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseByID - 2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseByID - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] - GetWarehouseByID - 4: %s", string(body))
		return nil, errors.New("failed to get warehouse by id")
	}

	var warehouseResponse WarehouseServiceResponse
	if err := json.Unmarshal(body, &warehouseResponse); err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseByID - 5: %v", err)
		return nil, err
	}

	return &warehouseResponse.Data, nil
}

// GetWarehouseProductStock implements WarehouseClientInterface.
func (w *WarehouseCLient) GetWarehouseProductStock(ctx context.Context, warehouseID uint, productID uint) (*WarehouseProductStockResponse, error) {
	url := fmt.Sprintf("%s/api/v1/warehouse-products/%d/detail/%d", w.urlWarehouseService, warehouseID, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseProductStock - 1: %v", err)
		return nil, err
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseProductStock - 2: %v", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseProductStock - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] - GetWarehouseProductStock - 4: %s", string(body))
		return nil, errors.New("failed to get warehouse product stock")
	}

	var warehouseProductStockResponse WarehouseProductStockServiceResponse
	if err := json.Unmarshal(body, &warehouseProductStockResponse); err != nil {
		log.Errorf("[WarehouseClient] - GetWarehouseProductStock - 5: %v", err)
		return nil, err
	}

	return &warehouseProductStockResponse.Data, nil
}

type WarehouseResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
	Phone   string `json:"phone"`
}

type WarehouseServiceResponse struct {
	Message string            `json:"message"`
	Data    WarehouseResponse `json:"data"`
	Error   string            `json:"error,omitempty"`
}

type WarehouseProductStockResponse struct {
	ID          uint `json:"id"`
	ProductID   uint `json:"product_id"`
	Stock       int  `json:"stock"`
	WarehouseID uint `json:"warehouse_id"`
}

type WarehouseProductStockServiceResponse struct {
	Message string                        `json:"message"`
	Data    WarehouseProductStockResponse `json:"data"`
	Error   string                        `json:"error,omitempty"`
}

func NewWarehouseClient(cfg configs.Config) WarehouseClientInterface {
	return &WarehouseCLient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlWarehouseService: cfg.App.UrlWarehouseService,
	}
}
