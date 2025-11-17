package controller

import (
	"micro-warehouse/warehouse-service/controller/request"
	"micro-warehouse/warehouse-service/controller/response"
	"micro-warehouse/warehouse-service/model"
	"micro-warehouse/warehouse-service/pkg/conv"
	"micro-warehouse/warehouse-service/pkg/pagination"
	"micro-warehouse/warehouse-service/pkg/validator"
	"micro-warehouse/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseControllerInterface interface {
	CreateWarehouse(ctx *fiber.Ctx) error
	GetAllWarehouses(ctx *fiber.Ctx) error
	GetWarehouseByID(ctx *fiber.Ctx) error
	UpdateWarehouse(ctx *fiber.Ctx) error
	DeleteWarehouse(ctx *fiber.Ctx) error
}

type warehouseController struct {
	warehouseUsecase usecase.WarehouseUsecaseInterface
}

// CreateWarehouse implements WarehouseControllerInterface.
func (w *warehouseController) CreateWarehouse(ctx *fiber.Ctx) error {
	var req request.CreateWarehouseRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Warehouse{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	if err := w.warehouseUsecase.CreateWarehouse(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create warehouse",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Warehouse created successfully",
	})
}

// DeleteWarehouse implements WarehouseControllerInterface.
func (w *warehouseController) DeleteWarehouse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	if err := w.warehouseUsecase.DeleteWarehouse(ctx.Context(), warehouseID); err != nil {
		log.Errorf("[WarehouseController] DeleteWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse deleted successfully",
	})
}

// GetAllWarehouses implements WarehouseControllerInterface.
func (w *warehouseController) GetAllWarehouses(ctx *fiber.Ctx) error {
	var req request.GetAllWarehouseRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	warehouses, totalRecords, err := w.warehouseUsecase.GetAllWarehouses(ctx.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all warehouses",
		})
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(totalRecords))
	
	var warehouseResponse []response.WarehouseResponse
	for _, warehouse := range warehouses {
		warehouseResponse = append(warehouseResponse, response.WarehouseResponse{
			ID:        warehouse.ID,
			Name:      warehouse.Name,
			Address:   warehouse.Address,
			Phone:     warehouse.Phone,
			Photo:     warehouse.Photo,
			CountProduct: len(warehouse.WarehouseProducts),
		})
	}
		
	response := response.GetAllWarehouseResponse{
		Warehouses: warehouseResponse,
		Pagination: pagination,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
		"message": "Warehouses fetched successfully",
	})
}

// GetWarehouseByID implements WarehouseControllerInterface.
func (w *warehouseController) GetWarehouseByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	warehouse, err := w.warehouseUsecase.GetWarehouseByID(ctx.Context(), warehouseID)
	if err != nil {
		log.Errorf("[WarehouseController] GetWarehouseByID - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouse by ID",
		})
	}

	respWarehouses := response.DetailWarehouseResponse{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Address: warehouse.Address,
		Phone:   warehouse.Phone,
		Photo:   warehouse.Photo,
	}

	for _, wp := range warehouse.WarehouseProducts {
		respWarehouses.WarehouseProducts = append(respWarehouses.WarehouseProducts, response.WarehouseProductResponse{
			ID:          wp.ID,
			WarehouseID: wp.WarehouseID,
			ProductID:   wp.ProductID,
			Stock:       wp.Stock,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    respWarehouses,
		"message": "Warehouse fetched successfully",
	})
}

// UpdateWarehouse implements WarehouseControllerInterface.
func (w *warehouseController) UpdateWarehouse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	var req request.CreateWarehouseRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Warehouse{
		ID:      warehouseID,
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	if err := w.warehouseUsecase.UpdateWarehouse(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse updated successfully",
	})
}

func NewWarehouseController(warehouseUsecase usecase.WarehouseUsecaseInterface) WarehouseControllerInterface {
	return &warehouseController{
		warehouseUsecase: warehouseUsecase,
	}
}
