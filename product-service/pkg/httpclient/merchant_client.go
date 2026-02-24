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

type MerchantClient struct {
	urlMerchantService string
	httpClient         *http.Client
	config             configs.Config
}

type MerchantProductStockResponse struct {
	ProductID  uint `json:"product_id"`
	TotalStock int  `json:"total_stock"`
}

type MerchantProductStockServiceResponse struct {
	Message string                       `json:"message"`
	Data    MerchantProductStockResponse `json:"data"`
	Error   string                       `json:"error,omitempty"`
}

func (mc *MerchantClient) generateInternalToken() (string, error) {
	return jwt.GenerateInternalToken(mc.config)
}

func NewMerchantClient(cfg configs.Config) *MerchantClient {
	return &MerchantClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlMerchantService: cfg.App.UrlMerchantService,
		config:             cfg,
	}
}

func (mc *MerchantClient) GetProductStockAcrossMerchants(ctx context.Context, productID uint) (int, error) {
	url := fmt.Sprintf("%s/api/v1/merchant-products/%d/total-stock", mc.urlMerchantService, productID)

	token, err := mc.generateInternalToken()
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 1: %v", err)
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 2: %v", err)
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 3: %v", err)
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 4: %v", err)
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 5: %s", string(body))
		return 0, errors.New("failed to get product stock across merchants")
	}

	var stockResp MerchantProductStockServiceResponse
	if err := json.Unmarshal(body, &stockResp); err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 6: %v", err)
		return 0, err
	}

	return stockResp.Data.TotalStock, nil
}

func (mc *MerchantClient) DeleteAllProductMerchantProducts(ctx context.Context, productID uint) error {
	url := fmt.Sprintf("%s/api/v1/merchant-products/product/%d", mc.urlMerchantService, productID)

	token, err := mc.generateInternalToken()
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 1: %v", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 2: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 3: %v", err)
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 4: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 5: %s", string(body))
		return errors.New("failed to delete all product merchant products")
	}

	return nil
}