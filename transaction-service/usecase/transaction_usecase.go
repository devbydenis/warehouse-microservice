package usecase

import (
	"context"
	"fmt"
	"micro-warehouse/transaction-service/model"
	"micro-warehouse/transaction-service/pkg/httpclient"
	"micro-warehouse/transaction-service/pkg/rabbitmq"
	"micro-warehouse/transaction-service/repository"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type TransactionUsecaseInterface interface {
	GetDashboardStats(ctx context.Context, userID uint) (int64, int64, int64, error)                       // sorting response: total revenue, total transaction, product sold
	GetDashboardStatsByMerchant(ctx context.Context, userID, merchantID uint) (int64, int64, int64, error) // sorting response: total revenue, total transaction, product sold

	GetTransactions(ctx context.Context, page int, limit int, search, sortBy, sortOrder string, merchantID uint) ([]model.Transaction, int64, error)
	CreateTransaction(ctx context.Context, transaction model.Transaction) (int64, error)

	// Midtrans update status transaction
	UpdatePaymentStatus(ctx context.Context, orderID string, paymentStatus, paymentMethod, transactionID, fraudStatus string) error
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepositoryInterface
	merchantClient  httpclient.MerchantClientInterface
	rabbitMQService *rabbitmq.RabbitMQService
	productClient   httpclient.ProductClientInterface
	userClient      httpclient.UserClientInterface
}

// CreateTransaction implements [TransactionUsecaseInterface].
func (t *transactionUsecase) CreateTransaction(ctx context.Context, transaction model.Transaction) (int64, error) {
	if err := t.validateProductStocks(ctx, transaction); err != nil {
		log.Errorf("[TransactionUsecase] CreateTransaction - 1: %v", err)
		return 0, err
	}

	transactionID, err := t.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		log.Errorf("[TransactionUsecase] CreateTransaction - 2: %v", err)
		return 0, err
	}

	go func() {
		if err := t.publishStockReducedEvent(ctx, transaction); err != nil {
			log.Errorf("[TransactionUsecase] CreateTransaction CreateTransaction - 3: %v", err)
		}
	}()

	return transactionID, nil
}

// GetDashboardStats implements [TransactionUsecaseInterface].
func (t *transactionUsecase) GetDashboardStats(ctx context.Context, userID uint) (int64, int64, int64, error) {
	user, err := t.userClient.GetUserByID(ctx, userID)
	if err != nil {
		log.Errorf("[TransactionUsecase] GetDashboardStats - 1: %v", err)
		return 0, 0, 0, err
	}

	isManager := false
	if user.RoleName == "Manager" {
		isManager = true
	}

	if !isManager {
		return 0, 0, 0, fmt.Errorf("akses ditolak: user bukan manager")
	}

	totalRevenue, totalTransaction, totalProductSold, err := t.transactionRepo.GetDashboardStats(ctx)
	if err != nil {
		log.Errorf("[TransactionUsecase] GetDashboardStats - 2: %v", err)
		return 0, 0, 0, err
	}

	return totalRevenue, totalTransaction, totalProductSold, nil
}

// GetDashboardStatsByMerchant implements [TransactionUsecaseInterface].
func (t *transactionUsecase) GetDashboardStatsByMerchant(ctx context.Context, userID, merchantID uint) (int64, int64, int64, error) {
	user, err := t.userClient.GetUserByID(ctx, userID)
	if err != nil {
		log.Errorf("[TransactionUsecase] GetDashboardStatsByMerchant - 1: %v", err)
		return 0, 0, 0, err
	}

	isManager := false
	if user.RoleName == "Manager" {
		isManager = true
	}

	if !isManager {
		return 0, 0, 0, fmt.Errorf("akses ditolak: user bukan manager")
	}

	merchant, err := t.merchantClient.GetMerchantByID(ctx, merchantID)
	if err != nil {
		log.Errorf("[TransactionUsecase] GetDashboardStatsByMerchant - 2: %v", err)
		return 0, 0, 0, err
	}

	if merchant.KeeperID != userID {
		log.Errorf("[TransactionUsecase] GetDashboardStatsByMerchant - 3: %v", err)
		return 0, 0, 0, fmt.Errorf("user tidak memiliki akses ke merchant")
	}

	totalRevenue, totalTransaction, totalProductSold, err := t.transactionRepo.GetDashboardStatsByMerchant(ctx, merchantID)
	if err != nil {
		log.Errorf("[TransactionUsecase] GetDashboardStatsByMerchant - 4: %v", err)
		return 0, 0, 0, err
	}

	return totalRevenue, totalTransaction, totalProductSold, nil

}

// GetTransactions implements [TransactionUsecaseInterface].
func (t *transactionUsecase) GetTransactions(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string, merchantID uint) ([]model.Transaction, int64, error) {
	transactions, total, err := t.transactionRepo.GetTransactions(ctx, page, limit, search, sortBy, sortOrder, merchantID)
	if err != nil {
		log.Errorf("[TransactionUsecase] GetTransactions - 1: %v", err)
		return nil, 0, err
	}

	for i := range transactions {
		if err := t.enrichTransactionWithProductData(ctx, &transactions[i]); err != nil {
			log.Warnf("[TransactionUsecase] GetTransactions - Failed to enrich transaction %d with product data: %v", transactions[i].ID, err)
			return nil, 0, err
		}

		if err := t.enrichTransactionWithMerchantData(ctx, &transactions[i]); err != nil {
			log.Errorf("[TransactionUsecase] GetTransactions - Failed to enrich transaction %d with merchant data: %v", err)
			return nil, 0, err
		}
	}

	return transactions, total, nil
}

// UpdatePaymentStatus implements [TransactionUsecaseInterface].
func (t *transactionUsecase) UpdatePaymentStatus(ctx context.Context, orderID string, paymentStatus string, paymentMethod string, transactionID string, fraudStatus string) error {
	return t.transactionRepo.UpdatePaymentStatus(ctx, orderID, paymentStatus, paymentMethod, transactionID, fraudStatus)
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepositoryInterface,
	merchantClient httpclient.MerchantClientInterface,
	rabbitMQService *rabbitmq.RabbitMQService,
	productClient httpclient.ProductClientInterface,
	userClient httpclient.UserClientInterface,
) TransactionUsecaseInterface {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		merchantClient:  merchantClient,
		rabbitMQService: rabbitMQService,
		productClient:   productClient,
		userClient:      userClient,
	}
}

func (tu *transactionUsecase) validateProductStocks(ctx context.Context, transaction model.Transaction) error {
	for _, product := range transaction.TransactionProducts {
		// Get stock information from merchant service
		merchantProduct, err := tu.merchantClient.GetMerchantProductStock(ctx, transaction.MerchantID, product.ProductID)
		if err != nil {
			log.Errorf("[TransactionUsecase] validateProductStocks - 3: %v", err)
			return err
		}

		// check if available stock is sufficent
		if merchantProduct.Stock < int(product.Quantity) {
			log.Errorf("[TransactionUsecase] validateProductStocks - 4: insufficient stock for product %d. Required: %d, Available: %d", product.ProductID, product.Quantity, merchantProduct.Stock)
			return fmt.Errorf("stock tidak mencukupi untuk product '%s'. Dibutuhkan: %d, Tersedia: %d",
				merchantProduct.ProductName, product.Quantity, merchantProduct.Stock)
		}

		log.Infof("[TransactionUsecase] validateProductStocks - Stock validation passed for product %d (%s). Required: %d, Available: %d", product.ProductID, merchantProduct.ProductName, product.Quantity, merchantProduct.Stock)
	}

	return nil
}

func (tu *transactionUsecase) publishStockReducedEvent(ctx context.Context, transaction model.Transaction) error {
	// prepare product for event
	var products []rabbitmq.StockReducedEventProduct
	for _, product := range transaction.TransactionProducts {
		products = append(products, rabbitmq.StockReducedEventProduct{
			ProductID: product.ProductID,
			Quantity:  int(product.Quantity),
		})
	}

	// create event
	event := rabbitmq.StockReducedEvent{
		MerchantID: transaction.MerchantID,
		Products:   products,
		OrderID:    transaction.OrderID,
		Timestamp:  time.Now(),
	}

	//publish event
	return tu.rabbitMQService.PublishStockReducedEvent(ctx, event)
}

func (tu *transactionUsecase) enrichTransactionWithProductData(ctx context.Context, transaction *model.Transaction) error {
	var products []httpclient.ProductResponse
	for _, tp := range transaction.TransactionProducts {
		product, err := tu.productClient.GetProductByID(ctx, tp.ProductID)
		if err != nil {
			log.Errorf("[TransactionUsecase] enrichTransactionWithProductData - 1: %v", err)
			return err
		}
		products = append(products, *product)
	}

	productMap := make(map[uint]httpclient.ProductResponse)
	for _, product := range products {
		productMap[product.ID] = product
	}

	for i := range transaction.TransactionProducts {
		tp := &transaction.TransactionProducts[i]
		if product, exists := productMap[tp.ProductID]; exists {
			tp.ProductName = product.Name
			tp.ProductPhoto = product.Thumbnail
			tp.ProductAbout = product.About
			tp.ProductCategoryID = product.Category.ID
			tp.ProductCategoryName = product.Category.Name
			tp.ProductCategoryPhoto = product.Category.Photo
		}
	}

	return nil
}

func (tu *transactionUsecase) enrichTransactionWithMerchantData(ctx context.Context, transaction *model.Transaction) error {
	merchant, err := tu.merchantClient.GetMerchantByID(ctx, transaction.MerchantID)
	if err != nil {
		log.Errorf("[TransactionUsecase] enrichTransactionWithMerchantData - 1: %v", err)
		return err
	}

	transaction.MerchantName = merchant.Name

	return nil
}
