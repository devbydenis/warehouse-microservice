package repository

import (
	"context"
	"micro-warehouse/product-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetAllProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Product, int64, error)
	GetProductById(ctx context.Context, id uint) (*model.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}

type productRepository struct {
	db *gorm.DB
}

// CreateProduct implements ProductRepositoryInterface.
func (p *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	select {
	case <-ctx.Done():
		log.Errorf("[ProductRepository] CreateProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		return p.db.WithContext(ctx).Create(product).Error
	}
}

// DeleteProduct implements ProductRepositoryInterface.
func (p *productRepository) DeleteProduct(ctx context.Context, id uint) error {
	select {
	case <-ctx.Done():
		log.Errorf("[ProductRepository] DeleteProduct - 1: %v", ctx.Err())
		return ctx.Err()
	default:
		modelProduct := model.Product{}
		if err := p.db.WithContext(ctx).First(&modelProduct, id).Error; err != nil {
			log.Errorf("[ProductRepository] DeleteProduct - 2: %v", ctx.Err())
			return err
		}

		return p.db.WithContext(ctx).Delete(&modelProduct).Error
	}
}

// GetAllProducts implements ProductRepositoryInterface.
func (p *productRepository) GetAllProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Product, int64, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[ProductRepository] GetAllProducts - 1: %v", ctx.Err())
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
		query := p.db.Model(&model.Product{})

		if search != "" {
			query = query.Where("name ILIKE ? OR tagline ILIKE ? OR about ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}

		var totalRecords int64
		if err := query.Count(&totalRecords).Error; err != nil {
			log.Errorf("[ProductRepository] GetAllProducts - 2: Failed to count products: %v", err)
			return nil, 0, err
		}

		var modelProducts []model.Product
		if err := query.
			Order(sortBy + " " + sortOrder).
			WithContext(ctx).
			Preload("Category").
			Offset(offset).
			Limit(limit).
			Find(&modelProducts).Error; err != nil {
			log.Errorf("[ProductRepository] GetAllProducts - 3: %v", err)
			return nil, 0, err
		}

		return modelProducts, totalRecords, nil
	}
}

// GetProductByBarcode implements ProductRepositoryInterface.
func (p *productRepository) GetProductByBarcode(ctx context.Context, barcode string) (*model.Product, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[ProductRepository] GetProductByBarcode - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
			modelProduct := model.Product{}
			if err := p.db.WithContext(ctx).Where("barcode = ?", barcode).Preload("Category").First(&modelProduct).Error; err != nil {
				log.Errorf("[ProductRepository] GetProductByBarcode - 2: %v", err)
				return nil, err
			}

			return &modelProduct, nil
	}
}

// GetProductById implements ProductRepositoryInterface.
func (p *productRepository) GetProductById(ctx context.Context, id uint) (*model.Product, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[ProductRepository] GetProductByID - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
			modelProduct := model.Product{}
			if err := p.db.WithContext(ctx).Where("id = ?", id).Preload("Category").First(&modelProduct).Error; err != nil {
				log.Errorf("[ProductRepository] GetProductByID - 2: %v", err)
				return nil, err
			}

			return &modelProduct, nil
	}
}

// UpdateProduct implements ProductRepositoryInterface.
func (p *productRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	select {
		case <-ctx.Done():
			log.Errorf("[ProductRepository] UpdateProduct - 1: %v", ctx.Err())
			return ctx.Err()
		default:
			existingProduct := model.Product{}
			if err := p.db.WithContext(ctx).Where("id = ?", product.ID).First(&existingProduct).Error; err != nil {
				log.Errorf("[ProductRepository] UpdateProduct - 2: %v", err)
				return err
			}

			updates := map[string]interface{}{
				"name": product.Name,
				"barcode": product.Barcode,
				"price": product.Price,
				"about": product.About,
				"category_id": product.CategoryID,
				"thumbnail": product.Thumbnail,
				"is_popular": product.IsPopular,
				"updated_at": product.UpdatedAt,
			}

			return p.db.WithContext(ctx).Model(&existingProduct).Updates(updates).Error
	}
}

func NewProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &productRepository{
		db: db,
	}
}
