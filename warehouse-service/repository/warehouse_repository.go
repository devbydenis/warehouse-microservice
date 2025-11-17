package repository

import (
	"context"
	"errors"
	"micro-warehouse/warehouse-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseRepositoryInterface interface {
	GetAllWarehouses(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Warehouse, int64, error)
	GetWarehouseByID(ctx context.Context, id uint) (*model.Warehouse, error)
	CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error
	DeleteWarehouse(ctx context.Context, id uint) error
}

type warehouseRepository struct {
	db *gorm.DB
}

// CreateWarehouse implements WarehouseRepositoryInterface.
func (w *warehouseRepository) CreateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] CreateWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return w.db.WithContext(ctx).Create(&warehouse).Error
	}
}

// DeleteWarehouse implements WarehouseRepositoryInterface.
func (w *warehouseRepository) DeleteWarehouse(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] DeleteWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelWarehouse := model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] DeleteWarehouse - 2: %v", err)
			return err
		}

		if len(modelWarehouse.WarehouseProducts) > 0 {
			log.Errorf("[WarehouseRepository] DeleteWarehouse - 3: %v", errors.New("Warehouse has products"))
			return errors.New("Warehouse has products")
		}

		return w.db.WithContext(ctx).Delete(&modelWarehouse).Error
	}
}

// GetAllWarehouses implements WarehouseRepositoryInterface.
func (w *warehouseRepository) GetAllWarehouses(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Warehouse, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetAllWarehouses - 1: %v", ctx.Err())
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

		// calculate offset
		offset := (page - 1) * limit

		// build query
		query := w.db.Model(&model.Warehouse{})

		if search != "" {
			query = query.Where("name ILIKE ? OR address ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		var totalRecords int64
		if err := query.Count(&totalRecords).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouses - 2: %v", err)
			return nil, 0, err
		}

		var modelWarehouses []model.Warehouse
		if err := query.
			Order(sortBy + " " + sortOrder).
			WithContext(ctx).
			Preload("WarehouseProducts").
			Offset(offset).
			Limit(limit).
			Find(&modelWarehouses).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetAllWarehouses - 3: %v", err)
			return nil, 0, err
		}

		return modelWarehouses, totalRecords, nil
	}
}

// GetWarehouseByID implements WarehouseRepositoryInterface.
func (w *warehouseRepository) GetWarehouseByID(ctx context.Context, id uint) (*model.Warehouse, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] GetWarehouseByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		var modelWarehouse model.Warehouse
		if err := w.db.WithContext(ctx).Where("id = ?", id).Preload("WarehouseProducts").First(&modelWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] GetWarehouseByID - 2: %v", err)
			return nil, err
		}

		return &modelWarehouse, nil
	}
}

// UpdateWarehouse implements WarehouseRepositoryInterface.
func (w *warehouseRepository) UpdateWarehouse(ctx context.Context, warehouse *model.Warehouse) error {
	select {
	case <-ctx.Done():
		log.Errorf("[WarehouseRepository] UpdateWarehouse - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingWarehouse := model.Warehouse{}
		if err := w.db.WithContext(ctx).Where("id = ?", warehouse.ID).First(&existingWarehouse).Error; err != nil {
			log.Errorf("[WarehouseRepository] UpdateWarehouse - 2: %v", err)
			return err
		}

		// update fields
		existingWarehouse.Name = warehouse.Name
		existingWarehouse.Address = warehouse.Address
		existingWarehouse.Phone = warehouse.Phone
		existingWarehouse.Photo = warehouse.Photo

		return w.db.WithContext(ctx).Save(&existingWarehouse).Error
	}
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepositoryInterface {
	return &warehouseRepository{
		db: db,
	}
}
