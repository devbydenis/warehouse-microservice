package repository

import (
	"context"
	"errors"
	"micro-warehouse/merchant-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// CRUD,
// getMerchantByProductIDAndMerchant,
// deleteAllProductMerchantProducts,
// getProductTotalStock,
// reduceStock,

type MerchantProductRepositoryInterface interface {
	CreateMerchantProduct(ctx context.Context, merchantProduct *model.MerchantProduct) error
	GetMerchantProductByID(ctx context.Context, id uint) (*model.MerchantProduct, error)
	GetMerchantProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string, merchantID, productID uint) ([]model.MerchantProduct, int64, error)
	GetMerchantProductByProductIDAndMerchantID(ctx context.Context, productID, merchantID uint) (*model.MerchantProduct, error)
	UpdateMerchantProduct(ctx context.Context, merchantProduct *model.MerchantProduct) error
	DeleteMerchantProduct(ctx context.Context, id uint) error
	DeleteAllProductMerchantProducts(ctx context.Context, productID uint) error

	GetProductTotalStock(ctx context.Context, productID uint) (int, error)
	ReduceStock(ctx context.Context, merchantID uint, productID uint, quantity int64) error
}

type merchantProductRepository struct {
	db *gorm.DB
}

// CreateMerchantProduct implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) CreateMerchantProduct(ctx context.Context, merchantProduct *model.MerchantProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] CreateMerchantProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return m.db.Create(merchantProduct).Error
	}
}

// DeleteAllProductMerchantProducts implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) DeleteAllProductMerchantProducts(ctx context.Context, productID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] DeleteAllProductMerchantProducts - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return m.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&model.MerchantProduct{}).Error
	}
}

// DeleteMerchantProduct implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) DeleteMerchantProduct(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] DeleteMerchantProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelMerchantProduct := model.MerchantProduct{}
		if err := m.db.WithContext(ctx).Where("id = ?", id).First(&modelMerchantProduct).Error; err != nil {
			log.Errorf("[MerchantProductRepository] DeleteMerchantProduct - 2: %v", err)
			return err
		}
		return m.db.WithContext(ctx).Delete(&modelMerchantProduct).Error
	}
}

// GetMerchantProductByID implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) GetMerchantProductByID(ctx context.Context, id uint) (*model.MerchantProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] GetMerchantProductByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var merchantProduct model.MerchantProduct
		if err := m.db.WithContext(ctx).Where("id = ?", id).First(&merchantProduct).Error; err != nil {
			log.Errorf("[MerchantProductRepository] GetMerchantProductByID - 2: %v", err)
			return nil, err
		}
		return &merchantProduct, nil
	}
}

// GetMerchantProductByProductIDAndMerchantID implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) GetMerchantProductByProductIDAndMerchantID(ctx context.Context, productID uint, merchantID uint) (*model.MerchantProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] GetMerchantProductByProductIDAndMerchantID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var merchantProduct model.MerchantProduct
		if err := m.db.WithContext(ctx).Where("product_id = ? AND merchant_id = ?", productID, merchantID).First(&merchantProduct).Error; err != nil {
			log.Errorf("[MerchantProductRepository] GetMerchantProductByProductIDAndMerchantID - 2: %v", err)
			return nil, err
		}
		return &merchantProduct, nil
	}
}

// GetMerchantProducts implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) GetMerchantProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string, merchantID uint, productID uint) ([]model.MerchantProduct, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] GetMerchantProducts - 1: %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		if sortBy == "" {
			sortBy = "created_at"
		}
		if sortOrder == "" {
			sortOrder = "desc"
		}

		var totalRecords int64
		modelMerchantProducts := []model.MerchantProduct{}

		offset := (page - 1) * limit
		query := m.db.WithContext(ctx).Model(&model.MerchantProduct{})

		if search != "" {
			searchTerm := "%" + search + "%"
			query = query.Where("stock::text ILIKE ?", searchTerm)
		}

		if merchantID != 0 {
			query = query.Where("merchant_id = ?", merchantID)
		}

		if productID != 0 {
			query = query.Where("product_id = ?", productID)
		}

		if err := query.Count(&totalRecords).Error; err != nil {
			log.Errorf("[MerchantProductRepository] GetMerchantProducts - 2: %v", err)
			return nil, 0, err
		}

		if err := query.Order(sortBy + " " + sortOrder).
			WithContext(ctx).
			Preload("Merchant").
			Offset(offset).
			Limit(limit).
			Find(&modelMerchantProducts).Error; err != nil {
			log.Errorf("[MerchantProductRepository] GetMerchantProducts - 3: %v", err)
			return nil, 0, err
		}
		return modelMerchantProducts, totalRecords, nil
	}
}

// GetProductTotalStock implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) GetProductTotalStock(ctx context.Context, productID uint) (int, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] GetProductTotalStock - 1: %v", ctx.Err())
		return 0, ctx.Err()
	default:
		var totalStock int
		if err := m.db.WithContext(ctx).Model(&model.MerchantProduct{}).
			Where("product_id = ?", productID).
			Select("SUM(stock)").
			Scan(&totalStock).Error; err != nil {
			log.Errorf("[MerchantProductRepository] GetProductTotalStock - 2: %v", err)
			return 0, err
		}
		return totalStock, nil
	}
}

// ReduceStock implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) ReduceStock(ctx context.Context, merchantID uint, productID uint, quantity int64) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] ReduceStock - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		var merchantProduct model.MerchantProduct
		err := m.db.WithContext(ctx).Where("merchant_id = ? AND product_id = ?", merchantID, productID).First(&merchantProduct).Error
		if err != nil {
			log.Errorf("[MerchantProductRepository] ReduceStock - 2: %v", err)
			return err
		}

		if merchantProduct.Stock < int(quantity) {
			log.Errorf("[MerchantProductRepository] ReduceStock - 3:", errors.New("stock is not enough"))
			return errors.New("stock is not enough")
		}

		newStock := merchantProduct.Stock - int(quantity)

		return m.db.WithContext(ctx).Model(&merchantProduct).Update("stock", newStock).Error
	}
}

// UpdateMerchantProduct implements MerchantProductRepositoryInterface.
func (m *merchantProductRepository) UpdateMerchantProduct(ctx context.Context, merchantProduct *model.MerchantProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantProductRepository] UpdateMerchantProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingMerchantProduct := model.MerchantProduct{}
		if err := m.db.WithContext(ctx).Where("id = ?", merchantProduct.ID).First(&existingMerchantProduct).Error; err != nil {
			log.Errorf("[MerchantProductRepository] UpdateMerchantProduct - 2: %v", err)
			return err
		}
		
		existingMerchantProduct.Stock = merchantProduct.Stock
		existingMerchantProduct.MerchantID = merchantProduct.MerchantID
		existingMerchantProduct.ProductID = merchantProduct.ProductID
		existingMerchantProduct.WarehouseID = merchantProduct.WarehouseID

		return m.db.WithContext(ctx).Save(&existingMerchantProduct).Error

	}
}

func NewMerchantProductRepository(db *gorm.DB) MerchantProductRepositoryInterface {
	return &merchantProductRepository{db: db}
}
