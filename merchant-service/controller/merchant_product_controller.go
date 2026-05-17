package controller

import (
	"micro-warehouse/merchant-service/controller/request"
	"micro-warehouse/merchant-service/controller/response"
	"micro-warehouse/merchant-service/model"
	"micro-warehouse/merchant-service/pkg/conv"
	"micro-warehouse/merchant-service/pkg/httpclient"
	"micro-warehouse/merchant-service/pkg/pagination"
	"micro-warehouse/merchant-service/pkg/validator"
	"micro-warehouse/merchant-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MerchantProductControllerInterface interface {
	CreateMerchantProduct(c *fiber.Ctx) error
	GetMerchantProductByID(c *fiber.Ctx) error
	GetMerchantProducts(c *fiber.Ctx) error
	GetMerchantProductByBarcode(c *fiber.Ctx) error
	UpdateMerchantProduct(c *fiber.Ctx) error
	DeleteMerchantProduct(c *fiber.Ctx) error
	DeleteAllProductMerchantProducts(c *fiber.Ctx) error
	GetProductTotalStock(c *fiber.Ctx) error
}

type merchantProductController struct {
	merchantProductUsecase usecase.MerchantProductUsecaseInterface
}

// CreateMerchantProduct implements MerchantProductControllerInterface.
// @Summary Create Merchant Product
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param product body request.CreateMerchantProductRequest true "Product"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products [post]
func (m *merchantProductController) CreateMerchantProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.CreateMerchantProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[MerchantProductController] CreateMerchantProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[MerchantProductController] CreateMerchantProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.MerchantProduct{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		Stock:       req.Stock,
		MerchantID:  req.MerchantID,
	}

	if err := m.merchantProductUsecase.CreateMerchantProduct(ctx, &reqModel); err != nil {
		log.Errorf("[MerchantProductController] CreateMerchantProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create merchant product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant product created successfully",
	})
}

// DeleteAllProductMerchantProducts implements MerchantProductControllerInterface.
// @Summary Delete All Product Merchant Products
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param product_id path uint true "Product ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products/product/{product_id} [delete]
func (m *merchantProductController) DeleteAllProductMerchantProducts(c *fiber.Ctx) error {
	ctx := c.Context()
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	if err := m.merchantProductUsecase.DeleteAllProductMerchantProducts(ctx, productIDUint); err != nil {
		log.Errorf("[MerchantProductController] DeleteAllProductMerchantProducts - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete all product merchant products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All product merchant products deleted successfully",
	})
}

// DeleteMerchantProduct implements MerchantProductControllerInterface.
// @Summary Delete Merchant Product
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param merchant_product_id path uint true "Merchant Product ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products/{merchant_product_id} [delete]
func (m *merchantProductController) DeleteMerchantProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	merchantProductID := c.Params("merchant_product_id")
	merchantProductIDUint := conv.StringToUint(merchantProductID)

	if err := m.merchantProductUsecase.DeleteMerchantProduct(ctx, merchantProductIDUint); err != nil {
		log.Errorf("[MerchantProductController] DeleteMerchantProduct - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete merchant product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant product deleted successfully",
	})
}

// GetMerchantProductByBarcode implements MerchantProductControllerInterface.
// @Summary Get Merchant Product By Barcode
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param barcode path string true "Barcode"
// @Param merchant_id query uint true "Merchant ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products/barcode/{barcode} [get]
func (m *merchantProductController) GetMerchantProductByBarcode(c *fiber.Ctx) error {
	ctx := c.Context()
	barcode := c.Params("barcode")
	merchantID := c.Params("merchant_id")
	merchantIDUint := conv.StringToUint(merchantID)
	if merchantIDUint == 0 {
		merchantIDUint = conv.StringToUint(c.Query("merchant_id"))
	}

	merchantProduct, product, warehouse, err := m.merchantProductUsecase.GetMerchantProductByBarcode(ctx, barcode, merchantIDUint)
	if err != nil {
		log.Errorf("[MerchantProductController] GetMerchantProductByBarcode - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get merchant product by barcode",
		})
	}

	productResponse := httpclient.MapProductResponseToMerchantProduct(product)
	warehouseResponse := httpclient.MapWarehouseResponseToMerchantProduct(warehouse)

	productResponse.ID = merchantProduct.ID
	productResponse.MerchantID = merchantProduct.MerchantID
	productResponse.ProductID = merchantProduct.ProductID
	productResponse.Stock = merchantProduct.Stock
	productResponse.WarehouseID = merchantProduct.WarehouseID
	productResponse.WarehouseName = warehouseResponse.WarehouseName
	productResponse.WarehousePhoto = warehouseResponse.WarehousePhoto
	productResponse.WarehousePhone = warehouseResponse.WarehousePhone

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant product fetched successfully",
		"data":    productResponse,
	})
}

// GetMerchantProductByID implements MerchantProductControllerInterface.
// @Summary Get Merchant Product By ID
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param merchant_product_id path uint true "Merchant Product ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products/{merchant_product_id} [get]
func (m *merchantProductController) GetMerchantProductByID(c *fiber.Ctx) error {
	ctx := c.Context()
	merchantProductID := c.Params("merchant_product_id")
	merchantProductIDUint := conv.StringToUint(merchantProductID)

	merchantProduct, product, warehouse, err := m.merchantProductUsecase.GetMerchantProductByID(ctx, merchantProductIDUint)
	if err != nil {
		log.Errorf("[MerchantProductController] GetMerchantProductByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get merchant product by ID",
		})
	}

	productResponse := httpclient.MapProductResponseToMerchantProduct(product)
	warehouseResponse := httpclient.MapWarehouseResponseToMerchantProduct(warehouse)

	productResponse.ID = merchantProduct.ID
	productResponse.MerchantID = merchantProduct.MerchantID
	productResponse.ProductID = merchantProduct.ProductID
	productResponse.Stock = merchantProduct.Stock
	productResponse.WarehouseID = merchantProduct.WarehouseID
	productResponse.WarehouseName = warehouseResponse.WarehouseName
	productResponse.WarehousePhoto = warehouseResponse.WarehousePhoto
	productResponse.WarehousePhone = warehouseResponse.WarehousePhone

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant product fetched successfully",
		"data":    productResponse,
	})
}

// GetMerchantProducts implements MerchantProductControllerInterface.
// @Summary Get Merchant Products
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param 			page query int false "Page"
// @Param 			limit query int false "Limit"
// @Param 			search query string false "Search"
// @Param 			sort_by query string false "Sort By"
// @Param 			sort_order query string false "Sort Order"
// @Param 			merchant_id query uint false "Merchant ID"
// @Param 			product_id query uint false "Product ID"
// @Param 			keeper_id query uint false "Keeper ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products [get]
func (m *merchantProductController) GetMerchantProducts(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetMerchantProductRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[MerchantProductController] GetMerchantProducts - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	merchantProducts, products, warehouses, total, err := m.merchantProductUsecase.GetMerchantProducts(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder, req.MerchantID, req.ProductID)
	if err != nil {
		log.Errorf("[MerchantProductController] GetMerchantProducts - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get merchant products",
		})
	}

	var productResponses []response.MerchantProduct
	productMap := make(map[uint]*httpclient.ProductResponse)
	warehouseMap := make(map[uint]*httpclient.WarehouseResponse)

	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for i := range warehouses {
		warehouseMap[warehouses[i].ID] = &warehouses[i]
	}

	for _, mp := range merchantProducts {
		productResponse := response.MerchantProduct{
			ID:          mp.ID,
			MerchantID:  mp.MerchantID,
			ProductID:   mp.ProductID,
			Stock:       mp.Stock,
			WarehouseID: mp.WarehouseID,
		}

		if product, exist := productMap[mp.ProductID]; exist {
			productResponse.ProductName = product.Name
			productResponse.ProductAbout = product.About
			productResponse.ProductPhoto = product.Thumbnail
			productResponse.ProductPrice = int(product.Price)
			productResponse.ProductCategory = product.Category.Name
			productResponse.ProductCategoryPhoto = product.Category.Photo
		}

		if warehouse, exist := warehouseMap[mp.WarehouseID]; exist {
			productResponse.WarehouseName = warehouse.Name
			productResponse.WarehousePhoto = warehouse.Photo
			productResponse.WarehousePhone = warehouse.Phone
		}

		productResponses = append(productResponses, productResponse)
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(total))
	response := response.GetAllMerchantProductsResponse{
		MerchantProducts: productResponses,
		Pagination:       pagination,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
		"message": "Merchant products fetched successfully",
	})
}

// GetProductTotalStock implements MerchantProductControllerInterface.
// @Summary Get Product Total Stock
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param product_id path uint true "Product ID"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products/{product_id}/total-stock [get]
func (m *merchantProductController) GetProductTotalStock(c *fiber.Ctx) error {
	ctx := c.Context()
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	totalStock, err := m.merchantProductUsecase.GetProductTotalStock(ctx, productIDUint)
	if err != nil {
		log.Errorf("[MerchantProductController] GetProductTotalStock - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get product total stock",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product total stock fetched successfully",
		"data":    totalStock,
	})
}

// UpdateMerchantProduct implements MerchantProductControllerInterface.
// @Summary Update Merchant Product
// @Tags Merchant Product
// @Accept json
// @Produce json
// @Param merchant_product_id path uint true "Merchant Product ID"
// @Param request body request.CreateMerchantProductRequest true "Request Body"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/merchant-products/{merchant_product_id} [put]
func (m *merchantProductController) UpdateMerchantProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	merchantProductID := c.Params("merchant_product_id")
	merchantProductIDUint := conv.StringToUint(merchantProductID)

	var req request.CreateMerchantProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[MerchantProductController] UpdateMerchantProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[MerchantProductController] UpdateMerchantProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.MerchantProduct{
		ID: merchantProductIDUint,
		ProductID: req.ProductID,
		WarehouseID: req.WarehouseID,
		Stock: req.Stock,
		MerchantID: req.MerchantID,
	}

	if err := m.merchantProductUsecase.UpdateMerchantProduct(ctx, &reqModel); err != nil {
		log.Errorf("[MerchantProductController] UpdateMerchantProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update merchant product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant product updated successfully",
	})
}

func NewMerchantProductController(merchantProductUsecase usecase.MerchantProductUsecaseInterface) MerchantProductControllerInterface {
	return &merchantProductController{
		merchantProductUsecase: merchantProductUsecase,
	}
}
