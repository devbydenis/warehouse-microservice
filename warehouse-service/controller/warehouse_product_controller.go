package controller

import (
	"micro-warehouse/warehouse-service/controller/request"
	"micro-warehouse/warehouse-service/controller/response"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/pkg/conv"
	"micro-warehouse/warehouse-service/pkg/httpclient"
	"micro-warehouse/warehouse-service/pkg/validator"
	"micro-warehouse/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseProductControllerInterface interface {
	GetDetailWarehouse(c *fiber.Ctx) error
	GetDetailWarehouseProductByID(c *fiber.Ctx) error
	CreateWarehouseProduct(c *fiber.Ctx) error
	GetWarehouseProductByWarehouseIDAndProductID(c *fiber.Ctx) error
	UpdateWarehouseProduct(c *fiber.Ctx) error
	DeleteWarehouseProduct(c *fiber.Ctx) error
	DeleteAllWarehouseProductByProductID(c *fiber.Ctx) error
	GetWarehouseProductByProductID(c *fiber.Ctx) error
	GetProductTotalStock(c *fiber.Ctx) error
}

type warehouseProductController struct {
	warehouseProductUsecase usecase.WarehouseProductUsecaseInterface
}

// CreateWarehouseProduct implements WarehouseProductController.
// @Summary Create warehouse product
// @Description Create warehouse product
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param warehouse_id path int true "Warehouse ID"
// @Param request body request.CreateWarehouseProductRequest true "Request body"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/{warehouse_id} [post]
// @Security Bearer
func (w *warehouseProductController) CreateWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.CreateWarehouseProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouse - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouse - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)

	reqModel := model.WarehouseProduct{
		WarehouseID: warehouseIDUint,
		ProductID:   req.ProductID,
		Stock:       req.Stock,
	}

	if err := w.warehouseProductUsecase.CreateWarehouseProduct(ctx, &reqModel); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create warehouse product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Warehouse product created successfully",
	})
}

// DeleteAllWarehouseProductByProductID implements WarehouseProductController.
// @Summary Delete all warehouse products by product ID
// @Description Delete all warehouse products by product ID
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/detail/products/{product_id} [delete]
// @Security Bearer
func (w *warehouseProductController) DeleteAllWarehouseProductByProductID(c *fiber.Ctx) error {
	ctx := c.Context()
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	if err := w.warehouseProductUsecase.DeleteAllWarehouseProductByProductID(ctx, productIDUint); err != nil {
		log.Errorf("[WarehouseProductController] DeleteAllWarehouseProductByProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete all warehouse products by product ID",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All warehouse products by product ID deleted successfully",
	})
}

// DeleteWarehouseProduct implements WarehouseProductController.
// @Summary Delete warehouse product
// @Description Delete warehouse product
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param warehouse_product_id path int true "Warehouse Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/detail/{warehouse_product_id} [delete]
// @Security Bearer
func (w *warehouseProductController) DeleteWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseProductID := c.Params("warehouse_product_id")
	warehouseProductIDUint := conv.StringToUint(warehouseProductID)

	if err := w.warehouseProductUsecase.DeleteWarehouseProduct(ctx, warehouseProductIDUint); err != nil {
		log.Errorf("[WarehouseProductController] DeleteWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete warehouse product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse product deleted successfully",
	})
}

// GetDetailWarehouse implements WarehouseProductController.
// @Summary Get detail warehouse
// @Description Get detail warehouse
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param warehouse_id path int true "Warehouse ID"
// @Success 200 {object} response.DetailWarehouseResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/{warehouse_id} [get]
// @Security Bearer
func (w *warehouseProductController) GetDetailWarehouse(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)

	warehouse, products, err := w.warehouseProductUsecase.GetDetailWarehouse(ctx, warehouseIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetDetailWarehouse - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get detail warehouse",
		})
	}

	respWarehouseProducts := response.DetailWarehouseResponse{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Address: warehouse.Address,
		Photo:   warehouse.Photo,
		Phone:   warehouse.Phone,
	}

	productMap := make(map[uint]*httpclient.ProductResponse)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for _, wp := range warehouse.WarehouseProducts {
		warehouseProduct := response.WarehouseProductResponse{
			ID:          wp.ID,
			WarehouseID: wp.WarehouseID,
			ProductID:   wp.ProductID,
			Stock:       wp.Stock,
		}

		if product, exist := productMap[wp.ProductID]; exist {
			warehouseProduct.ProductName = product.Name
			warehouseProduct.ProductAbout = product.About
			warehouseProduct.ProductPhoto = product.Thumbnail
			warehouseProduct.ProductPrice = int(product.Price)
			warehouseProduct.ProductCategory = product.Category.Name
			warehouseProduct.ProductCategoryPhoto = product.Category.Photo
		}

		respWarehouseProducts.WarehouseProducts = append(respWarehouseProducts.WarehouseProducts, warehouseProduct)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    respWarehouseProducts,
		"message": "Detail warehouse fetched successfully",
	})
}

// GetDetailWarehouseProductByID implements WarehouseProductController.
// @Summary Get detail warehouse product by ID
// @Description Get detail warehouse product by ID
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param warehouse_product_id path int true "Warehouse Product ID"
// @Success 200 {object} response.GetDetailWarehouseProductByIDResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/detail/products/{warehouse_product_id}/warehouses [get]
// @Security Bearer
func (w *warehouseProductController) GetDetailWarehouseProductByID(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseProductID := c.Params("warehouse_product_id")
	warehouseProductIDUint := conv.StringToUint(warehouseProductID)

	warehouseProduct, product, err := w.warehouseProductUsecase.GetDetailWarehouseProductByID(ctx, warehouseProductIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetDetailWarehouseProductByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get detail warehouse product by ID",
		})
	}

	respWarehouseProduct := response.GetDetailWarehouseProductByIDResponse{
		ID:             warehouseProduct.ID,
		WarehouseID:    warehouseProduct.WarehouseID,
		ProductID:      warehouseProduct.ProductID,
		Stock:          warehouseProduct.Stock,
		WarehouseName:  warehouseProduct.Warehouse.Name,
		WarehousePhoto: warehouseProduct.Warehouse.Photo,
		WarehousePhone: warehouseProduct.Warehouse.Phone,
		ProductName:    product.Name,
		ProductBarcode: product.Barcode,
		ProductPrice:   int(product.Price),
		ProductAbout:   product.About,
		ProductThumbnail: product.Thumbnail,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    respWarehouseProduct,
		"message": "Detail warehouse product by ID fetched successfully",
	})
}

// GetProductTotalStock implements WarehouseProductController.
// @Summary Get product total stock
// @Description Get product total stock
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param product_id path int true "Product ID"
// @Success 200 {object} response.ProductTotalStockResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/detail/products/{product_id}/total-stock [get]
// @Security Bearer
func (w *warehouseProductController) GetProductTotalStock(c *fiber.Ctx) error {
	ctx := c.Context()
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	totalStock, err := w.warehouseProductUsecase.GetProductTotalStock(ctx, productIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetProductTotalStock - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get product total stock",
		})
	}

	resp := response.ProductTotalStockResponse{
		ProductID:  productIDUint,
		TotalStock: totalStock,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": 	resp,
		"message": "Product total stock fetched successfully",
	})
}

// GetWarehouseProductByProductID implements WarehouseProductController.
// @Summary Get warehouse product by product ID
// @Description Get warehouse product by product ID
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param product_id path int true "Product ID"
// @Success 200 {object} []response.WarehouseResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/detail/products/{product_id} [get]
// @Security Bearer
func (w *warehouseProductController) GetWarehouseProductByProductID(c *fiber.Ctx) error {
	ctx := c.Context()
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	warehouseProducts, err := w.warehouseProductUsecase.GetWarehouseProductByProductID(ctx, productIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetWarehouseProductByProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouse products by product ID",
		})
	}

	resp := []response.WarehouseResponse{}
	for _, wp := range warehouseProducts {
		resp = append(resp, response.WarehouseResponse{
			ID:        wp.Warehouse.ID,
			Name:      wp.Warehouse.Name,
			Address:   wp.Warehouse.Address,
			Phone:     wp.Warehouse.Phone,
			Photo:     wp.Warehouse.Photo,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse products by product ID fetched successfully",
	})
}

// GetWarehouseProductByWarehouseIDAndProductID implements WarehouseProductController.
// @Summary Get warehouse product by warehouse ID and product ID
// @Description Get warehouse product by warehouse ID and product ID
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param warehouse_id path int true "Warehouse ID"
// @Param product_id path int true "Product ID"
// @Success 200 {object} response.WarehouseProductResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/{warehouse_id}/detail/{product_id} [get]
// @Security Bearer
func (w *warehouseProductController) GetWarehouseProductByWarehouseIDAndProductID(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	warehouseProduct, err := w.warehouseProductUsecase.GetWarehouseProductByWarehouseIDAndProductID(ctx, warehouseIDUint, productIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetWarehouseProductByWarehouseIDAndProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouse product by warehouse ID and product ID",
		})
	}

	resp := response.WarehouseProductResponse{
		ID:          warehouseProduct.ID,
		WarehouseID: warehouseProduct.WarehouseID,
		ProductID:   warehouseProduct.ProductID,
		Stock:       warehouseProduct.Stock,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Warehouse product by warehouse ID and product ID fetched successfully",
	})
}

// UpdateWarehouseProduct implements WarehouseProductController.
// @Summary Update warehouse product
// @Description Update warehouse product
// @Tags Warehouse Products
// @Accept  json
// @Produce  json
// @Param warehouse_product_id path int true "Warehouse Product ID"
// @Param warehouse_id path int true "Warehouse ID"
// @Param request body request.CreateWarehouseProductRequest true "Request body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/warehouse-products/detail/{warehouse_product_id}/{warehouse_id} [put]
// @Security Bearer
func (w *warehouseProductController) UpdateWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseProductID := c.Params("warehouse_product_id")
	warehouseProductIDUint := conv.StringToUint(warehouseProductID)
	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)
	
	var req request.CreateWarehouseProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.WarehouseProduct{
		ID:          warehouseProductIDUint,
		WarehouseID: warehouseIDUint,
		ProductID:   req.ProductID,
		Stock:       req.Stock,
	}

	if err := w.warehouseProductUsecase.UpdateWarehouseProduct(ctx, &reqModel); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update warehouse product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse product updated successfully",
	})
}

func NewWarehouseProductController(warehouseProductUsecase usecase.WarehouseProductUsecaseInterface) WarehouseProductControllerInterface {
	return &warehouseProductController{
		warehouseProductUsecase: warehouseProductUsecase,
	}
}
