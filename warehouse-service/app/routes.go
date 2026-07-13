package app

import (
	"github.com/gofiber/fiber/v2"
	
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "micro-warehouse/warehouse-service/docs"
)

func SetupRoutes(app *fiber.App, c* Container) {
	api := app.Group("/api/v1")

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	warehouses := api.Group("/warehouses")
	warehouses.Post("/", c.WarehouseController.CreateWarehouse)
	warehouses.Get("/", c.WarehouseController.GetAllWarehouses)
	warehouses.Get("/:id", c.WarehouseController.GetWarehouseByID)
	warehouses.Put("/:id", c.WarehouseController.UpdateWarehouse)
	warehouses.Delete("/:id", c.WarehouseController.DeleteWarehouse)

	warehouseProducts := api.Group("/warehouse-products")
	warehouseProducts.Post("/:warehouse_id", c.WarehouseProductController.CreateWarehouseProduct)
	warehouseProducts.Get("/:warehouse_id", c.WarehouseProductController.GetDetailWarehouse)
	warehouseProducts.Get("/:warehouse_id/detail/:product_id", c.WarehouseProductController.GetWarehouseProductByWarehouseIDAndProductID)
	warehouseProducts.Put("/detail/:warehouse_product_id/:warehouse_id", c.WarehouseProductController.UpdateWarehouseProduct)
	warehouseProducts.Delete("/detail/:warehouse_product_id", c.WarehouseProductController.DeleteWarehouseProduct)
	warehouseProducts.Delete("/detail/products/:product_id", c.WarehouseProductController.DeleteAllWarehouseProductByProductID)
	warehouseProducts.Get("/detail/products/:product_id/total-stock", c.WarehouseProductController.GetProductTotalStock)
	warehouseProducts.Get("/detail/products/:product_id", c.WarehouseProductController.GetWarehouseProductByProductID)
	warehouseProducts.Get("/detail/products/:warehouse_product_id/warehouses", c.WarehouseProductController.GetDetailWarehouseProductByID)

	api.Post("/upload-warehouse", c.UploadController.UploadPhoto)
}