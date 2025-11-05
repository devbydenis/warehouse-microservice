package controller

import (
	"micro-warehouse/product-service/controller/request"
	"micro-warehouse/product-service/controller/response"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/pkg/conv"
	"micro-warehouse/product-service/pkg/pagination"
	"micro-warehouse/product-service/pkg/validator"
	"micro-warehouse/product-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ProductControllerInterface interface {
	CreateProduct(ctx *fiber.Ctx) error
	GetAllProducts(ctx *fiber.Ctx) error
	GetProductByID(ctx *fiber.Ctx) error
	GetProductByBarcode(ctx *fiber.Ctx) error
	UpdateProduct(ctx *fiber.Ctx) error
	DeleteProduct(ctx *fiber.Ctx) error
}

type productController struct {
	productUsecase usecase.ProductUsecaseInterface
}

// CreateProduct implements ProductControllerInterface.
func (p *productController) CreateProduct(ctx *fiber.Ctx) error {
	var req request.CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[ProductController] CreateProduct - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[ProductController] CreateProduct - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Product{
		Name: req.Name,
		Barcode: req.Barcode,
		CategoryID: req.CategoryID,
		Thumbnail: req.Thumbnail,
		About: req.About,
		Price: req.Price,
		IsPopular: req.IsPopular,
	}

	if err := p.productUsecase.CreateProduct(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[ProductController] CreateProduct - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create product",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
	})
}

// DeleteProduct implements ProductControllerInterface.
func (p *productController) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idUint := conv.StringToUint(id)

	if err := p.productUsecase.DeleteProduct(ctx.Context(), idUint); err != nil {
		log.Errorf("[ProductController] DeleteProduct - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete product",
		})
	}
	
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// GetAllProducts implements ProductControllerInterface.
func (p *productController) GetAllProducts(ctx *fiber.Ctx) error {
	var req request.GetAllProductRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[ProductController] GetAllProducts - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	products, total, err := p.productUsecase.GetAllProducts(ctx.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[ProductController] GetAllProducts - 2: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get products",
		})
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(total))

	var productsResponse []response.ProductResponse
	for _, product := range products {
		productsResponse = append(productsResponse, response.ProductResponse{
			ID: product.ID,
			Name: product.Name,
			Barcode: product.Barcode,
			CategoryID: product.CategoryID,
			Thumbnail: product.Thumbnail,
			About: product.About,
			Price: int(product.Price),
			IsPopular: product.IsPopular,
		})
	}

	response := response.GetAllProductResponse{
		Products: productsResponse,
		Pagination: pagination,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products fetched successfully",
		"data": response,
	})
}

// GetProductByBarcode implements ProductControllerInterface.
func (p *productController) GetProductByBarcode(ctx *fiber.Ctx) error {
	barcode := ctx.Params("barcode")

	product, err := p.productUsecase.GetProductByBarcode(ctx.Context(), barcode)
	if err != nil {
		log.Errorf("[ProductController] GetProductByBarcode - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get product",
		})
	}

	response := response.ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Barcode: product.Barcode,
		CategoryID: product.CategoryID,
		Thumbnail: product.Thumbnail,
		About: product.About,
		Price: int(product.Price),
		IsPopular: product.IsPopular,
		Category: response.CategoryResponse{
			ID: product.Category.ID,
			Name: product.Category.Name,
			Tagline: product.Category.Tagline,
			Photo: product.Category.Photo,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product fetched successfully",
		"data": response,
	})
}

// GetProductByID implements ProductControllerInterface.
func (p *productController) GetProductByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idUint := conv.StringToUint(id)

	product, err := p.productUsecase.GetProductById(ctx.Context(), idUint)
	if err != nil {
		log.Errorf("[ProductController] GetProductByID - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get product",
		})
	}

	response := response.ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Barcode: product.Barcode,
		CategoryID: product.CategoryID,
		Thumbnail: product.Thumbnail,
		About: product.About,
		Price: int(product.Price),
		IsPopular: product.IsPopular,
		Category: response.CategoryResponse{
			ID: product.Category.ID,
			Name: product.Category.Name,
			Tagline: product.Category.Tagline,
			Photo: product.Category.Photo,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product fetched successfully",
		"data": response,
	})
}

// UpdateProduct implements ProductControllerInterface.
func (p *productController) UpdateProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idUint := conv.StringToUint(id)

	var req request.CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[ProductController] UpdateProduct - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[ProductController] UpdateProduct - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Product{
		ID: idUint,
		Name: req.Name,
		Barcode: req.Barcode,
		CategoryID: req.CategoryID,
		Thumbnail: req.Thumbnail,
		About: req.About,
		Price: req.Price,
		IsPopular: req.IsPopular,
	}

	if err := p.productUsecase.UpdateProduct(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[ProductController] UpdateProduct - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update product",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

func NewProductController(productUsecase usecase.ProductUsecaseInterface) ProductControllerInterface {
	return &productController{
		productUsecase: productUsecase,
	}
}
