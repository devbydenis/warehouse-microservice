package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-warehouse/transaction-service/configs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type MerchantClientInterface interface {
	GetMerchantsByKeeperID(ctx context.Context, keeperID uint) ([]Merchant, error)
	GetMerchantByID(ctx context.Context, merchantID uint) (*Merchant, error)
	GetMerchantProducts(ctx context.Context, merchantID uint) ([]MerchantProduct, error)
	GetMerchantProductStock(ctx context.Context, merchantID uint, productID uint) (*MerchantProduct, error)
}

type MerchantCLient struct {
	urlMerchantService string
	httpClient         *http.Client
}

// GetMerchantByID implements [MerchantClientInterface].
func (m *MerchantCLient) GetMerchantByID(ctx context.Context, merchantID uint) (*Merchant, error) {
	url := fmt.Sprintf("%s/api/v1/merchants/%d", m.urlMerchantService, merchantID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantByID - 1: %v", err)
		return nil, err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantByID - 2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantByID - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] GetMerchantByID - 4: %s", string(body))
		return nil, errors.New("failed to get merchant by id")
	}

	var response struct {
		Data Merchant `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Errorf("[MerchantCLient] GetMerchantByID - 5: %v", err)
		return nil, err
	}

	return &response.Data, nil
}

// GetMerchantProductStock implements [MerchantClientInterface].
func (m *MerchantCLient) GetMerchantProductStock(ctx context.Context, merchantID uint, productID uint) (*MerchantProduct, error) {
	url := fmt.Sprintf("%s/api/v1/merchant-products?merchant_id=%d&product_id=%d", m.urlMerchantService, merchantID, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantProductStock - 1: %v", err)
		return nil, err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantProductStock - 2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantProductStock - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] GetMerchantProductStock - 2: %s", string(body))
		return nil, errors.New("failed to get merchant product stock")
	}

	var response struct {
		Data struct {
			MerchantProducts []MerchantProduct `json:"merchant_products"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Errorf("[MerchantCLient] GetMerchantProductStock - 5: %v", err)
		return nil, err
	}

	for _, product := range response.Data.MerchantProducts {
		if product.ProductID == productID {
			log.Infof("[MerchantCLient] GetMerchantProductStock - Found product %d with stock %d", productID, product.Stock)
			return &product, nil
		}
	}

	return nil, errors.New("product not found")
}

// GetMerchantProducts implements [MerchantClientInterface].
func (m *MerchantCLient) GetMerchantProducts(ctx context.Context, merchantID uint) ([]MerchantProduct, error) {
	url := fmt.Sprintf("%s/api/v1/merchant-products?merchant_id=%d", m.urlMerchantService, merchantID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantProducts - 1: %v", err)
		return nil, err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantProducts - 2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantProducts - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] GetMerchantProducts - 2: %s", string(body))
		return nil, errors.New("failed to get merchant product stock")
	}

	var response struct {
		Data []MerchantProduct `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Errorf("[MerchantCLient] GetMerchantProducts - 5: %v", err)
		return nil, err
	}

	return response.Data, nil
}

// GetMerchantsByKeeperID implements [MerchantClientInterface].
func (m *MerchantCLient) GetMerchantsByKeeperID(ctx context.Context, keeperID uint) ([]Merchant, error) {
	url := fmt.Sprintf("%s/api/v1/merchants?keeper_id=%d", m.urlMerchantService, keeperID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantsByKeeperID - 1: %v", err)
		return nil, err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantsByKeeperID - 2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] GetMerchantsByKeeperID - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] GetMerchantsByKeeperID - 2: %s", string(body))
		return nil, errors.New("failed to get merchant product stock")
	}

	var response struct {
		Data []Merchant `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Errorf("[MerchantCLient] GetMerchantsByKeeperID - 5: %v", err)
		return nil, err
	}

	return response.Data, nil
}

type Merchant struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	KeeperID uint   `json:"keeper_id"`
}

type MerchantProduct struct {
	ID                   uint   `json:"id"`
	MerchantID           uint   `json:"merchant_id"`
	ProductID            uint   `json:"product_id"`
	ProductName          string `json:"product_name"`
	ProductAbout         string `json:"product_about"`
	ProductPhoto         string `json:"product_photo"`
	ProductPrice         uint   `json:"product_price"`
	ProductCategory      string `json:"product_category"`
	ProductCategoryPhoto string `json:"product_category_photo"`
	Stock                int    `json:"stock"`
	WarehouseID          uint   `json:"warehouse_id"`
	WarehouseName        string `json:"warehouse_name"`
	WarehousePhoto       string `json:"warehouse_photo"`
	WarehousePhone       string `json:"warehouse_phone"`
}

func NewMerchantClient(cfg configs.Config) MerchantClientInterface {
	return &MerchantCLient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlMerchantService: cfg.App.UrlMerchantService,
	}
}
