package controller

import (
	"fmt"
	"micro-warehouse/transaction-service/controller/request"
	"micro-warehouse/transaction-service/controller/response"
	"micro-warehouse/transaction-service/model"
	"micro-warehouse/transaction-service/pkg/conv"
	"micro-warehouse/transaction-service/pkg/midtrans"
	"micro-warehouse/transaction-service/pkg/pagination"
	"micro-warehouse/transaction-service/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TransactionControllerInterface interface {
	CreateTransaction(ctx *fiber.Ctx) error
	GetTransactions(c *fiber.Ctx) error
	MidtransCallback(c *fiber.Ctx) error

	GetManagerDashboard(c *fiber.Ctx) error
	GetDashboardByMerchant(c *fiber.Ctx) error
}

type transactionController struct {
	transactionUsecase usecase.TransactionUsecaseInterface
	midtransService    midtrans.MidtransServiceInterface
}

// CreateTransaction implements [TransactionControllerInterface].
// @Summary      Create transaction
// @Description  Create a new transaction
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Param      transaction body request.CreateTransactionWithProductRequest true "Transaction data"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/transactions [post]
// @Security Bearer
func (t *transactionController) CreateTransaction(ctx *fiber.Ctx) error {
	var req request.CreateTransactionWithProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[TransactionController] CreateTransaction - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	orderID := fmt.Sprintf("ORDER_%d_%d", time.Now().Unix(), req.MerchantID)

	var subtotal int64
	var items []midtrans.TransactionItem

	for _, product := range req.Products {
		productSubtotal := int64(product.Price) * product.Quantity
		subtotal += productSubtotal

		items = append(items, midtrans.TransactionItem{
			ID:       fmt.Sprintf("%d", product.ProductID),
			Price:    product.Price,
			Quantity: product.Quantity,
			Name:     fmt.Sprintf("Product %d", product.ProductID),
		})
	}

	taxTotal := int64(float64(subtotal) * 0.11)
	grandTotal := subtotal + taxTotal

	transaction := model.Transaction{
		Name:          req.Name,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		MerchantID:    req.MerchantID,
		SubTotal:      subtotal,
		TaxTotal:      taxTotal,
		Notes:         req.Notes,
		GrandTotal:    grandTotal,
		Currency:      "IDR",
		OrderID:       orderID,
		PaymentStatus: model.PaymentStatusPending,
		PaymentMethod: model.PaymentMethodQRIS,
	}

	for _, product := range req.Products {
		transaction.TransactionProducts = append(transaction.TransactionProducts, model.TransactionProduct{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
			Price:     product.Price,
			Subtotal:  product.Price * product.Quantity,
		})
	}

	_, err := t.transactionUsecase.CreateTransaction(ctx.Context(), transaction)
	if err != nil {
		log.Errorf("[TransactionController] CreateTransaction - 2: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create transaction",
		})
	}

	midtransReq := midtrans.CreateTransactionRequest{
		OrderID:       orderID,
		Amount:        int64(grandTotal),
		Items:         items,
		CustomerName:  req.Name,
		CustomerEmail: req.Email,
		CustomerPhone: req.Phone,
		Notes:         req.Notes,
	}

	midtransRes, err := t.midtransService.CreateTransaction(midtransReq)
	if err != nil {
		log.Errorf("[TransactionController] CreateTransaction - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create transaction",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaction created successfully",
		"data": fiber.Map{
			"payment_token": midtransRes.PaymentToken,
			"order_id":      midtransRes.OrderID,
		},
	})
}

// GetDashboardByMerchant implements [TransactionControllerInterface].
// @Summary      Get dashboard by merchant
// @Description  Get dashboard statistics by merchant
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Param      	 merchant_id path string true "Merchant ID"
// @Param      	 user_id query uint true "User ID"
// @Success      200  {object} response.DashboardByMerchantResponse
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/dashboard/keeper/merchant/{merchant_id} [get]
// @Security Bearer
func (t *transactionController) GetDashboardByMerchant(c *fiber.Ctx) error {
	ctx := c.Context()

	merchantIDStr := c.Params("merchant_id")
	merchantID := conv.StringToUint(merchantIDStr)

	userIDStr := c.Query("user_id")
	userID := conv.StringToUint(userIDStr)

	totalRevenue, totalTransactions, productsSold, err := t.transactionUsecase.GetDashboardStatsByMerchant(ctx, userID, merchantID)
	if err != nil {
		log.Errorf("[TransactionController] GetDashboardByMerchant - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get dashboard stats by merchant",
		})
	}

	response := response.DashboardByMerchantResponse{
		DashboardResponse: response.DashboardResponse{
			TotalRevenue:      totalRevenue,
			TotalTransactions: totalTransactions,
			ProductsSold:      productsSold,
		},
		Merchant: response.MerchantSummary{
			ID: merchantID,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Dashboard stats fetched successfully",
	})

}

// GetManagerDashboard implements [TransactionControllerInterface].
// @Summary      Get dashboard by manager
// @Description  Get dashboard statistics by manager
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Param      	 user_id query uint true "User ID"
// @Success      200  {object} response.DashboardResponse
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/dashboard/manager [get]
// @Security Bearer
func (t *transactionController) GetManagerDashboard(c *fiber.Ctx) error {
	ctx := c.Context()

	managerIDStr := c.Query("user_id")
	managerID := conv.StringToUint(managerIDStr)

	totalRevenue, totalTransactions, productsSold, err := t.transactionUsecase.GetDashboardStats(ctx, managerID)
	if err != nil {
		log.Errorf("[TransactionController] GetManagerDashboard - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get dashboard stats",
		})
	}

	response := response.DashboardResponse{
		TotalRevenue:      totalRevenue,
		TotalTransactions: totalTransactions,
		ProductsSold:      productsSold,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Dashboard stats fetched successfully",
	})
}

// GetTransactions implements [TransactionControllerInterface].
// @Summary      Get transactions
// @Description  Get all transactions
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param sort_by query string false "Sort By"
// @Param sort_order query string false "Sort Order"
// @Param merchant_id query string false "Merchant ID"
// @Success      200  {object} response.GetAllTransactionsResponse
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/transactions [get]
// @Security Bearer
func (t *transactionController) GetTransactions(c *fiber.Ctx) error {
	ctx := c.Context()

	query := request.GetAllTransactionRequest{}
	if err := c.QueryParser(&query); err != nil {
		log.Errorf("[TransactionController] GetTransactions - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query parameters",
		})
	}

	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}

	merchantID := conv.StringToUint(query.MerchantID)

	transactions, total, err := t.transactionUsecase.GetTransactions(ctx, query.Page, query.Limit, query.Search, query.SortBy, query.SortOrder, merchantID)
	if err != nil {
		log.Errorf("[TransactionController] GetTransactions - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get transactions",
		})
	}

	var transactionResponses []response.TransactionResponse
	for _, transaction := range transactions {
		var transactionProductResponses []response.TransactionProductResponse
		for _, tp := range transaction.TransactionProducts {
			transactionProductResponses = append(transactionProductResponses, response.TransactionProductResponse{
				ID:            tp.ID,
				ProductID:     tp.ProductID,
				ProductName:   tp.ProductName,
				ProductPhoto:  tp.ProductPhoto,
				ProductAbout:  tp.ProductAbout,
				Quantity:      tp.Quantity,
				Price:         tp.Price,
				SubTotal:      tp.Subtotal,
				TransactionID: tp.TransactionID,
				Category: struct {
					ID    uint   `json:"id"`
					Name  string `json:"name"`
					Photo string `json:"photo"`
				}{
					ID:    tp.ProductCategoryID,
					Name:  tp.ProductCategoryName,
					Photo: tp.ProductCategoryPhoto,
				},
			})
		}

		transactionResponses = append(transactionResponses, response.TransactionResponse{
			ID:                  transaction.ID,
			Name:                transaction.Name,
			Phone:               transaction.Phone,
			Email:               transaction.Email,
			Address:             transaction.Address,
			SubTotal:            transaction.SubTotal,
			TaxTotal:            transaction.TaxTotal,
			GrandTotal:          transaction.GrandTotal,
			MerchantID:          transaction.MerchantID,
			MerchantName:        transaction.MerchantName,
			PaymentStatus:       transaction.PaymentStatus,
			PaymentMethod:       transaction.PaymentMethod,
			TransactionCode:     transaction.TransactionCode,
			OrderID:             transaction.OrderID,
			Notes:               transaction.Notes,
			TransactionProducts: transactionProductResponses,
		})
	}

	paginationInfo := pagination.CalculatePagination(query.Page, query.Limit, int(total))

	response := response.GetAllTransactionsResponse{
		Transactions: transactionResponses,
		Pagination:   paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Transactions fetched successfully",
	})
}

// MidtransCallback implements [TransactionControllerInterface].
// @Summary      Midtrans callback
// @Description  Handle Midtrans callback
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Param        request body request.MidtransCallbackRequest true "Midtrans callback request"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/midtrans/callback [post]
// @Security Bearer
func (t *transactionController) MidtransCallback(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.MidtransCallbackRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[TransactionController] MidtransCallback - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := t.transactionUsecase.UpdatePaymentStatus(ctx, req.OrderID, req.TransactionStatus, req.PaymentType, req.TransactionID, req.FraudStatus); err != nil {
		log.Errorf("[TransactionController] MidtransCallback - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update payment status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment status updated successfully",
	})
}

func NewTransactionController(transactionUsecase usecase.TransactionUsecaseInterface, midtransService midtrans.MidtransServiceInterface) TransactionControllerInterface {
	return &transactionController{
		transactionUsecase: transactionUsecase,
		midtransService:    midtransService,
	}
}
