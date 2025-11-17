package repository

import (
	"context"
	"micro-warehouse/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseProductRepositoryInterface interface {
	GetDetailWarehouse(ctx context.Context, warehouseID uint) (*model.Warehouse, error)
	GetDetailWarehouseProductByID(ctx context.Context, warehouseProductID uint) (*model.WarehouseProduct, error)
	CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	GetWarehouseProductByWarehouseIDAndProductID(ctx context.Context, warehouseID, productID uint) (*model.WarehouseProduct, error)
	UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error
	DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error
	DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error
	GetWarehouseProductByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error)
	GetProductTotalStock(ctx context.Context, productID uint) (int, error)
}

type warehouseProductRepository struct {
	db *gorm.DB
}

// CreateWarehouseProduct implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) CreateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] CreateWarehouseProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(warehouseProduct).Error
	}
}

// DeleteAllWarehouseProductByProductID implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) DeleteAllWarehouseProductByProductID(ctx context.Context, productID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] DeleteAllWarehouseProductByProductID - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		err := w.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&model.WarehouseProduct{}).Error
		if err != nil {
			log.Errorf("[WarehouseProductRepository] DeleteAllWarehouseProductByProductID - 2: %v", err)
			return err
		}

		return nil
	}
}

// DeleteWarehouseProduct implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) DeleteWarehouseProduct(ctx context.Context, warehouseProductID uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] DeleteWarehouseProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouseProduct := &model.WarehouseProduct{}
		err := w.db.WithContext(ctx).Where("id = ?", warehouseProductID).First(&modelWarehouseProduct).Error
		if err != nil {
			log.Errorf("[WarehouseProductRepository] DeleteWarehouseProduct - 2: %v", err)
			return err
		}

		return w.db.WithContext(ctx).Delete(&modelWarehouseProduct).Error
	}
}

// GetDetailWarehouse implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) GetDetailWarehouse(ctx context.Context, warehouseID uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetDetailWarehouse - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouse model.Warehouse
		if err := w.db.
			WithContext(ctx).
			Where("id = ?", warehouseID).
			Select("id", "name", "photo", "phone").
			Preload("WarehouseProducts").
			First(&warehouse).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetDetailWarehouse - 2: %v", err)
			return nil, err
		}

		return &warehouse, nil
	}
}

// GetDetailWarehouseProductByID implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) GetDetailWarehouseProductByID(ctx context.Context, warehouseProductID uint) (*model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetDetailWarehouseProductByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouseProduct model.WarehouseProduct
		if err := w.db.
			WithContext(ctx).
			Where("id = ?", warehouseProductID).
			Select("id", "product_id", "stock", "warehouse_id").
			Preload("Warehouse").
			First(&warehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetDetailWarehouseProductByID - 2: %v", err)
			return nil, err
		}
		return &warehouseProduct, nil
	}
}

// GetProductTotalStock implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) GetProductTotalStock(ctx context.Context, productID uint) (int, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetProductTotalStock - 1: %v", ctx.Err())
		return 0, ctx.Err()
	default:
		var totalStock int
		if err := w.db.
			WithContext(ctx).
			Model(&model.WarehouseProduct{}).
			Where("product_id = ?", productID).
			Select("COALESCE(SUM(stock), 0)").
			Scan(&totalStock).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetProductTotalStock - 2: %v", err)
			return 0, err
		}
		return totalStock, nil
	}
}

// GetWarehouseProductByProductID implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) GetWarehouseProductByProductID(ctx context.Context, productID uint) ([]model.WarehouseProduct, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[WarehouseProductRepository] GetWarehouseProductByProductID - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
			var warehouseProducts []model.WarehouseProduct
			if err := w.db.
				WithContext(ctx).
				Where("product_id = ?", productID).
				Preload("Warehouse").
				Find(&warehouseProducts).Error; err != nil {
				log.Errorf("[WarehouseProductRepository] GetWarehouseProductByProductID - 2: %v", err)
				return nil, err
			}
			
			return warehouseProducts, nil
	}
}

// GetWarehouseProductByWarehouseIDAndProductID implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) GetWarehouseProductByWarehouseIDAndProductID(ctx context.Context, warehouseID uint, productID uint) (*model.WarehouseProduct, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] GetWarehouseProductByWarehouseIDAndProductID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var warehouseProduct model.WarehouseProduct
		if err := w.db.WithContext(ctx).
			Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
			First(&warehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] GetWarehouseProductByWarehouseIDAndProductID - 2: %v", err)
			return nil, err
		}
		
		return &warehouseProduct, nil
	}
}

// UpdateWarehouseProduct implements WarehouseProductRepositoryInterface.
func (w *warehouseProductRepository) UpdateWarehouseProduct(ctx context.Context, warehouseProduct *model.WarehouseProduct) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseProductRepository] UpdateWarehouseProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingWarehouseProduct := &model.WarehouseProduct{}
		if err := w.db.WithContext(ctx).Where("id = ?", warehouseProduct.ID).First(&existingWarehouseProduct).Error; err != nil {
			log.Errorf("[WarehouseProductRepository] UpdateWarehouseProduct - 2: %v", err)
			return err
		}

		existingWarehouseProduct.Stock = warehouseProduct.Stock
		existingWarehouseProduct.WarehouseID = warehouseProduct.WarehouseID
		existingWarehouseProduct.ProductID = warehouseProduct.ProductID

		return w.db.WithContext(ctx).Save(&existingWarehouseProduct).Error
	}
}

func NewWarehouseProductRepository(db *gorm.DB) WarehouseProductRepositoryInterface {
	return &warehouseProductRepository{
		db: db,
	}
}
