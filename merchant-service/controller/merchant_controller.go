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

type MerchantControllerInterface interface {
	CreateMerchant(c *fiber.Ctx) error
	GetAllMerchants(c *fiber.Ctx) error
	GetMerchantById(c *fiber.Ctx) error
	UpdateMerchant(c *fiber.Ctx) error
	DeleteMerchant(c *fiber.Ctx) error
}

type merchantController struct {
	merchantUsecase usecase.MerchantUsecaseInterface
}

// CreateMerchant implements MerchantControllerInterface.
func (m *merchantController) CreateMerchant(c *fiber.Ctx) error {
	var req request.CreateMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[MerchantController] CreateMerchant - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(&req); err != nil {
		log.Errorf("[MerchantController] CreateMerchant - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	reqModel := model.Merchant{
		Name:     req.Name,
		KeeperID: req.KeeperID,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	if err := m.merchantUsecase.CreateMerchant(c.Context(), &reqModel); err != nil {
		log.Errorf("[MerchantController] CreateMerchant - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create merchant",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Merchant created successfully",
	})
}

// DeleteMerchant implements MerchantControllerInterface.
func (m *merchantController) DeleteMerchant(c *fiber.Ctx) error {
	id := c.Params("id")
	merchantID := conv.StringToUint(id)

	if err := m.merchantUsecase.DeleteMerchant(c.Context(), merchantID); err != nil {
		log.Errorf("[MerchantController] DeleteMerchant - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete merchant",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant deleted successfully",
	})

}

// GetAllMerchants implements MerchantControllerInterface.
func (m *merchantController) GetAllMerchants(c *fiber.Ctx) error {
	var req request.GetMerchantProductRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[MerchantController] GetAllMerchants - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request query",
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.KeeperID != 0 {
		productMap := make(map[uint]*httpclient.ProductResponse)
		warehouseMap := make(map[uint]*httpclient.WarehouseResponse)

		merchant, products, warehouses, err := m.merchantUsecase.GetMerchantByKeeperID(c.Context(), req.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetAllMerchants - 2: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get merchants",
			})
		}

		if merchant.ID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No merchant found for this keeper",
			})
		}

		keeperNames, err := m.merchantUsecase.GetKeeperName(c.Context(), req.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetAllMerchants - 3: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get keeper name",
			})
		}

		merchantResponse := response.MerchantWithProductResponse{
			ID:         merchant.ID,
			Name:       merchant.Name,
			Address:    merchant.Address,
			Photo:      merchant.Photo,
			Phone:      merchant.Phone,
			KeeperID:   merchant.KeeperID,
			KeeperName: keeperNames,
		}

		if len(merchant.MerchantProducts) > 0 {
			for i := range products {
				productMap[products[i].ID] = &products[i]
			}

			for i := range warehouses {
				warehouseMap[warehouses[i].ID] = &warehouses[i]
			}

			for _, mp := range merchant.MerchantProducts {
				productResponse := response.MerchantProduct{
					ID:          mp.ID,
					MerchantID:  mp.MerchantID,
					ProductID:   mp.ProductID,
					WarehouseID: mp.WarehouseID,
					Stock:       mp.Stock,
				}

				if product, exists := productMap[mp.ProductID]; exists {
					productResponse.ProductName = product.Name
					productResponse.ProductAbout = product.About
					productResponse.ProductPhoto = product.Thumbnail
					productResponse.ProductPrice = int(product.Price)
					productResponse.ProductCategory = product.Category.Name
					productResponse.ProductCategoryPhoto = product.Category.Photo
				}

				if warehouse, exists := warehouseMap[mp.WarehouseID]; exists {
					productResponse.WarehouseName = warehouse.Name
					productResponse.WarehousePhoto = warehouse.Photo
					productResponse.WarehousePhone = warehouse.Phone

				}

				merchantResponse.MerchantProducts = append(merchantResponse.MerchantProducts, productResponse)
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Merchant found",
			"data":    merchantResponse,
		})
	}

	merchants, total, err := m.merchantUsecase.GetAllMerchants(c.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[MerchantController] GetAllMerchants - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get merchants",
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(total))
	if len(merchants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message":    "No merchants found",
			"data":       nil,
			"pagination": paginationInfo,
		})
	}

	var merchantResponse []response.MerchantResponse
	for _, merchant := range merchants {
		keeperName, err := m.merchantUsecase.GetKeeperName(c.Context(), merchant.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetAllMerchants - 3: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get keeper name",
			})
		}
		merchantResponse = append(merchantResponse, response.MerchantResponse{
			ID:           merchant.ID,
			Name:         merchant.Name,
			Address:      merchant.Address,
			Photo:        merchant.Photo,
			Phone:        merchant.Phone,
			KeeperID:     merchant.KeeperID,
			KeeperName:   keeperName,
			ProductCount: len(merchant.MerchantProducts),
		})
	}

	resp := response.MerchantPaginationResponse{
		Message:    "Merchants fetched successfully",
		Data:       merchantResponse,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       resp,
		"pagination": paginationInfo,
		"message":    "Merchants fetched successfully",
	})
}

// GetMerchantById implements MerchantControllerInterface.
func (m *merchantController) GetMerchantById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id := conv.StringToUint(idStr)

	merchant, err := m.merchantUsecase.GetMerchantByID(c.Context(), id)
	if err != nil {
		log.Errorf("[MerchantController] GetMerchantById - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get merchant",
		})
	}

	var keeperName string
	if merchant.KeeperID != 0 {
		keeperName, err = m.merchantUsecase.GetKeeperName(c.Context(), merchant.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetMerchantById - 2: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get keeper name",
			})
		}
	}

	response := response.MerchantResponse{
		ID:           merchant.ID,
		Name:         merchant.Name,
		Address:      merchant.Address,
		Photo:        merchant.Photo,
		Phone:        merchant.Phone,
		KeeperID:     merchant.KeeperID,
		KeeperName:   keeperName,
		ProductCount: len(merchant.MerchantProducts),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Merchant found",
	})
}

// UpdateMerchant implements MerchantControllerInterface.
func (m *merchantController) UpdateMerchant(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id := conv.StringToUint(idStr)

	var req request.CreateMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[MerchantController] UpdateMerchant - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[MerchantController] UpdateMerchant - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	reqModel := model.Merchant{
		ID:       id,
		Name:     req.Name,
		KeeperID: req.KeeperID,
		Address:  req.Address,
		Photo:    req.Photo,
		Phone:    req.Phone,
	}

	if err := m.merchantUsecase.UpdateMerchant(c.Context(), &reqModel); err != nil {
		log.Errorf("[MerchantController] UpdateMerchant - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update merchant",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Merchant updated successfully",
	})
}

func NewMerchantController(merchantUsecase usecase.MerchantUsecaseInterface) MerchantControllerInterface {
	return &merchantController{
		merchantUsecase: merchantUsecase,
	}
}
