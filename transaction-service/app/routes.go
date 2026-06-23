package app

import (
	"github.com/gofiber/fiber/v2"
	
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "micro-warehouse/transaction-service/docs"
	
)

func SetupRoutes(app *fiber.App, container *Container) {
	
	app.Post("/api/v1/midtrans/callback", container.TransactionController.MidtransCallback)
	
	api := app.Group("/api/v1")

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	dashboard := api.Group("/dashboard")
	dashboard.Get("/manager", container.TransactionController.GetManagerDashboard)
	dashboard.Get("/keeper/merchant/:merchant_id", container.TransactionController.GetDashboardByMerchant)

	transaction := api.Group("/transactions")
	transaction.Post("/", container.TransactionController.CreateTransaction)
	transaction.Get("/", container.TransactionController.GetTransactions)
}