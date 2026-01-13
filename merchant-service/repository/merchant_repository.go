package repository

import (
	"context"
	"micro-warehouse/merchant-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// create, get all (include pagination), get by id, update, delete
// get merchant by keeper id,
type MerchantRepositoryInterface interface {
	CreateMerchant(ctx context.Context, merchant *model.Merchant) error
	GetAllMerchants(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Merchant, int64, error)
	GetMerchantByID(ctx context.Context, id uint) (*model.Merchant, error)
	UpdateMerchant(ctx context.Context, merchant *model.Merchant) error
	DeleteMerchant(ctx context.Context, id uint) error
	GetMerchantByKeeperID(ctx context.Context, keeperID uint) (*model.Merchant, error)
}

type merchantRepository struct {
	db *gorm.DB
}

// CreateMerchant implements MerchantRepositoryInterface.
func (m *merchantRepository) CreateMerchant(ctx context.Context, merchant *model.Merchant) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantRepository] CreateMerchant - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return m.db.WithContext(ctx).Create(merchant).Error
	}
}

// DeleteMerchant implements MerchantRepositoryInterface.
func (m *merchantRepository) DeleteMerchant(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantRepository] DeleteMerchant - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return m.db.WithContext(ctx).Delete(&model.Merchant{}, id).Error
	}
}

// GetAllMerchants implements MerchantRepositoryInterface.
func (m *merchantRepository) GetAllMerchants(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Merchant, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantRepository] GetAllMerchants - 1: %v", ctx.Err())
		return nil, 0, ctx.Err()
	default:
		modelMerchants := []model.Merchant{}
		totalRecords := int64(0)

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

		query := m.db.WithContext(ctx).Model(&model.Merchant{})

		if search != "" {
			query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+search+"%", "%"+search+"%")
		}

		if err := query.Count(&totalRecords).Error; err != nil {
			log.Errorf("[MerchantRepository] GetAllMerchants - 2: %v", err)
			return nil, 0, err
		}

		if err := query.Preload("MerchantProducts").Order(sortBy + " " + sortOrder).WithContext(ctx).Offset(offset).Limit(limit).Find(&modelMerchants).Error; err != nil {
			log.Errorf("[MerchantRepository] GetAllMerchants - 3: %v", err)
			return nil, 0, err
		}

		return modelMerchants, totalRecords, nil
	}
}

// GetMerchantByID implements MerchantRepositoryInterface.
func (m *merchantRepository) GetMerchantByID(ctx context.Context, id uint) (*model.Merchant, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantRepository] GetMerchantByID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelMerchant := model.Merchant{}

		if err := m.db.WithContext(ctx).Where("id = ?", id).Preload("MerchantProducts").First(&modelMerchant).Error; err != nil {
			log.Errorf("[MerchantRepository] GetMerchantByID - 2: %v", err)
			return nil, err
		}

		return &modelMerchant, nil
	}
}

// GetMerchantByKeeperID implements MerchantRepositoryInterface.
func (m *merchantRepository) GetMerchantByKeeperID(ctx context.Context, keeperID uint) (*model.Merchant, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantRepository] GetMerchantByKeeperID - 1: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelMerchant := model.Merchant{}
		if err := m.db.WithContext(ctx).Where("keeper_id = ?", keeperID).Preload("MerchantProducts").First(&modelMerchant).Error; err != nil {
			log.Errorf("[MerchantREpository] GetMerchantByKeeperID - 2: %v", err)
			return nil, err
		}

		return &modelMerchant, nil
	}
}

// UpdateMerchant implements MerchantRepositoryInterface.
func (m *merchantRepository) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	select {
	case <-ctx.Done():
		log.Errorf("[MerchantRepository] UpdateMerchant - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		existingMerchant := model.Merchant{}
		if err := m.db.WithContext(ctx).Where("id = ?", merchant.ID).First(&existingMerchant).Error; err != nil {
			log.Errorf("[MerchantRepository] UpdateMerchant - 2: %v", err)
			return err
		}

		existingMerchant.Name = merchant.Name
		existingMerchant.Address = merchant.Address
		existingMerchant.Photo = merchant.Photo
		existingMerchant.Phone = merchant.Phone

		return m.db.WithContext(ctx).Save(&existingMerchant).Error
	}
}

func NewMerchantRepository(db *gorm.DB) MerchantRepositoryInterface {
	return &merchantRepository{
		db: db,
	}
}
