package app

import (
	_ "micro-warehouse/merchant-service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, c *Container) {
	api := app.Group("/api/v1")

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	merchants := api.Group("/merchants")
	merchants.Post("/", c.MerchantController.CreateMerchant)
	merchants.Get("/", c.MerchantController.GetAllMerchants)
	merchants.Get("/:id", c.MerchantController.GetMerchantById)
	merchants.Put("/:id", c.MerchantController.UpdateMerchant)
	merchants.Delete("/:id", c.MerchantController.DeleteMerchant)

	merchantProducts := api.Group("/merchant-products")
	merchantProducts.Post("/", c.MerchantProductController.CreateMerchantProduct)
	merchantProducts.Get("/:merchant_product_id", c.MerchantProductController.GetMerchantProductByID)
	merchantProducts.Get("/", c.MerchantProductController.GetMerchantProducts)
	merchantProducts.Get("/barcode/:barcode", c.MerchantProductController.GetMerchantProductByBarcode)
	merchantProducts.Put("/:merchant_product_id", c.MerchantProductController.UpdateMerchantProduct)
	merchantProducts.Delete("/:merchant_product_id", c.MerchantProductController.DeleteMerchantProduct)
	merchantProducts.Delete("/product/:product_id", c.MerchantProductController.DeleteAllProductMerchantProducts)
	merchantProducts.Get("/:product_id/total-stock", c.MerchantProductController.GetProductTotalStock)

	api.Post("/upload-merchant", c.UploadController.UploadMerchantPhoto)
}
