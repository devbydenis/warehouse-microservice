package usecase

import (
	"context"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/repository"
)

type ProductUsecaseInterface interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetAllProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Product, int64, error)
	GetProductById(ctx context.Context, id uint) (*model.Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id uint) error
}

type productUsecase struct {
	productRepo repository.ProductRepositoryInterface
}

// CreateProduct implements ProductUsecaseInterface.
func (p *productUsecase) CreateProduct(ctx context.Context, product *model.Product) error {
	return p.productRepo.CreateProduct(ctx, product)
}

// DeleteProduct implements ProductUsecaseInterface.
func (p *productUsecase) DeleteProduct(ctx context.Context, id uint) error {
	return p.productRepo.DeleteProduct(ctx, id)
}

// GetAllProducts implements ProductUsecaseInterface.
func (p *productUsecase) GetAllProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Product, int64, error) {
	return p.productRepo.GetAllProducts(ctx, page, limit, search, sortBy, sortOrder)
}

// GetProductByBarcode implements ProductUsecaseInterface.
func (p *productUsecase) GetProductByBarcode(ctx context.Context, barcode string) (*model.Product, error) {
	return p.productRepo.GetProductByBarcode(ctx, barcode)
}

// GetProductById implements ProductUsecaseInterface.
func (p *productUsecase) GetProductById(ctx context.Context, id uint) (*model.Product, error) {
	return p.productRepo.GetProductById(ctx, id)
}

// UpdateProduct implements ProductUsecaseInterface.
func (p *productUsecase) UpdateProduct(ctx context.Context, product *model.Product) error {
	return p.productRepo.UpdateProduct(ctx, product)
}

func NewProductUsecase(productRepo repository.ProductRepositoryInterface) ProductUsecaseInterface {
	return &productUsecase{
		productRepo: productRepo,
	}
}
