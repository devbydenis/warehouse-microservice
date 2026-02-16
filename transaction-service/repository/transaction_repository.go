package repository

import (
	"context"
	"errors"
	"micro-warehouse/transaction-service/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// get data overview dashboard manager/keeper
// create transaction
// update status transaction

type TransactionRepositoryInterface interface {
	GetDashboardStats(ctx context.Context) (int64, int64, int64, error)
	GetDashboardStatsByMerchant(ctx context.Context, merchantID uint) (int64, int64, int64, error)

	GetTransactions(ctx context.Context, page int, limit int, search, sortBy, sortOrder string, merchantID uint) ([]model.Transaction, int64, error)
	CreateTransaction(ctx context.Context, transaction model.Transaction) (int64, error)

	// Midtrans update status transaction
	UpdatePaymentStatus(ctx context.Context, orderID string, paymentStatus, paymentMethod, transactionID, fraudStatus string) error
}

type transactionRepository struct {
	db *gorm.DB
}

// CreateTransaction implements [TransactionRepositoryInterface].
func (t *transactionRepository) CreateTransaction(ctx context.Context, transaction model.Transaction) (int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[TransactionRepository] CreateTransaction - 1: %v", ctx.Err())
		return 0, ctx.Err()
	default:
		tx := t.db.WithContext(ctx).Begin()
		if tx.Error != nil {
			log.Errorf("[TransactionRepository] CreateTransaction - 2: %v", tx.Error)
			return 0, tx.Error
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Errorf("[TransactionRepository] CreateTransaction - 3: %v", r)

			}
		}()

		products := transaction.TransactionProducts
		transaction.TransactionProducts = nil

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			log.Errorf("[TransactionRepository] CreateTransaction - 4: %v", err)
			return 0, err
		}

		for _, product := range products {
			modelTransactionProduct := model.TransactionProduct{
				ProductID:     product.ProductID,
				Quantity:      product.Quantity,
				Price:         product.Price,
				Subtotal:      product.Subtotal,
				TransactionID: transaction.ID,
			}

			if err := tx.Create(&modelTransactionProduct).Error; err != nil {
				tx.Rollback()
				log.Errorf("[TransactionRepository] CreateTransaction - 5: %v", err)
				return 0, err
			}
		}
		if err := tx.Commit().Error; err != nil {
			log.Errorf("[TransactionRepository] CreateTransaction - 6: %v", err)
			return 0, err
		}
		return int64(transaction.ID), nil
	}
}

// GetDashboardStats implements [TransactionRepositoryInterface].
func (t *transactionRepository) GetDashboardStats(ctx context.Context) (int64, int64, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[TransactionRepository] GetDashboardStats - 1: %v", ctx.Err())
		return 0, 0, 0, ctx.Err()
	default:
		var totalRevenue int64
		var totalTransactions int64
		var productsSold int64

		var result struct {
			TotalRevenue      int64 `json:"total_revenue"`
			TotalTransactions int64 `json:"total_transactions"`
			ProductsSold      int64 `json:"products_sold"`
		}

		err := t.db.WithContext(ctx).Model(&model.Transaction{}).
			Where("payment_status = ?", model.PaymentStatusSuccess).
			Select("COALESCE(SUM(grand_total), 0) as total_revenue, COUNT(*) as total_transactions").
			Scan(&result).Error

		if err != nil {
			log.Errorf("[TransactionRepository] GetDashboardStats - 2: %v", err)
			return 0, 0, 0, err
		}

		totalRevenue = result.TotalRevenue
		totalTransactions = result.TotalTransactions

		err = t.db.WithContext(ctx).Model(&model.TransactionProduct{}).
			Joins("JOIN transactions ON transaction_products.transaction_id = transactions.id").
			Where("transactions.payment_status = ?", model.PaymentStatusSuccess).
			Select("COALESCE(SUM(transaction_products.quantity), 0) as products_sold").
			Scan(&productsSold).Error

		if err != nil {
			log.Errorf("[TransactionRepository] GetDashboardStats - 3: %v", err)
			return 0, 0, 0, err
		}

		return totalRevenue, totalTransactions, productsSold, nil
	}
}

// GetDashboardStatsByMerchant implements [TransactionRepositoryInterface].
func (t *transactionRepository) GetDashboardStatsByMerchant(ctx context.Context, merchantID uint) (int64, int64, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[TransactionRepository] GetDashboardStatsByMerchant - 1: %v", ctx.Err())
		return 0, 0, 0, ctx.Err()
	default:

		if merchantID == 0 {
			return 0, 0, 0, fiber.ErrBadRequest
		}

		var totalRevenue int64
		var totalTransactions int64
		var productsSold int64

		var result struct {
			TotalRevenue      int64 `json:"total_revenue"`
			TotalTransactions int64 `json:"total_transactions"`
			ProductsSold      int64 `json:"products_sold"`
		}

		err := t.db.WithContext(ctx).Model(&model.Transaction{}).
			Where("merchant_id = ? AND payment_status = ?", merchantID, model.PaymentStatusSuccess).
			Select("COALESCE(SUM(grand_total), 0) as total_revenue, COUNT(*) as total_transactions").
			Scan(&result).Error

		if err != nil {
			log.Errorf("[TransactionRepository] GetDashboardStatsByMerchant - 2: %v", err)
			return 0, 0, 0, err
		}

		totalRevenue = result.TotalRevenue
		totalTransactions = result.TotalTransactions

		err = t.db.WithContext(ctx).Model(&model.TransactionProduct{}).
			Joins("JOIN transactions ON transaction_products.transaction_id = transactions.id").
			Where("transactions.merchant_id = ? AND transactions.payment_status = ?", merchantID, model.PaymentStatusSuccess).
			Select("COALESCE(SUM(transaction_products.quantity), 0) as products_sold").
			Scan(&productsSold).Error

		if err != nil {
			log.Errorf("[TransactionRepository] GetDashboardStatsByMerchant - 3: %v", err)
			return 0, 0, 0, err
		}

		return totalRevenue, totalTransactions, productsSold, nil
	}
}

// GetTransactions implements [TransactionRepositoryInterface].
func (t *transactionRepository) GetTransactions(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string, merchantID uint) ([]model.Transaction, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[TransactionRepository] GetTransactions - 1: %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		if sortBy == "" {
			sortBy = "created_at"
		}
		if sortOrder == "" {
			sortOrder = "desc"
		}

		offset := (page - 1) * limit

		baseSql := t.db.WithContext(ctx).Model(&model.Transaction{}).
			Preload("TransactionProducts")

		if search != "" {
			searchTerm := "%" + search + "%"
			baseSql = baseSql.Where("name ILIKE ? OR phone ILIKE ?", searchTerm, searchTerm)
		}

		if merchantID != 0 {
			baseSql = baseSql.Where("merchant_id = ?", merchantID)
		}

		var totalRecords int64
		if err := baseSql.Count(&totalRecords).Error; err != nil {
			log.Errorf("[TransactionRepository] GetAllTransactions - Count: %v", errors.New("failed to count transactions"))
			return nil, 0, err
		}

		var transactions []model.Transaction
		err := baseSql.WithContext(ctx).
			Preload("TransactionProducts").
			Order(sortBy + " " + sortOrder).
			Limit(limit).
			Offset(offset).
			Find(&transactions).Error

		if err != nil {
			log.Errorf("[TransactionRepository] GetTransactions - 3: %v", err)
			return nil, 0, err
		}

		return transactions, totalRecords, nil
	}
}

// UpdatePaymentStatus implements [TransactionRepositoryInterface].
func (t *transactionRepository) UpdatePaymentStatus(ctx context.Context, orderID string, paymentStatus string, paymentMethod string, transactionID string, fraudStatus string) error {
	select {
	case <-ctx.Done():
		log.Errorf("[TransactionRepository] UpdatePaymentStatus - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		
		if err := t.db.WithContext(ctx).Model(&model.Transaction{}).Where("order_id = ?", orderID).First(&model.Transaction{}).Error; err != nil {
			log.Errorf("[TransactionRepository] UpdatePaymentStatus - 2: %v", err)
			return err
		}

		updates := map[string]interface{}{
			"payment_status": paymentStatus,
		}

		if paymentMethod != "" {
			updates["payment_method"] = paymentMethod
		}

		if transactionID != "" {
			updates["transaction_code"] = transactionID
		}

		if fraudStatus != "" {
			updates["fraud_status"] = fraudStatus
		}

		err := t.db.WithContext(ctx).Model(&model.Transaction{}).
			Where("order_id = ?", orderID).
			Updates(updates).Error

		if err != nil {
			log.Errorf("[TransactionRepository] UpdatePaymentStatus - 3: %v", err)
			return err
		}

		return nil
	}
}

func NewTransactionRepository(db *gorm.DB) TransactionRepositoryInterface {
	return &transactionRepository{db: db}
}
