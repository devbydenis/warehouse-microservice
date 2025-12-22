package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, c *Container) {
	api := app.Group("/api/v1")

	merchants := api.Group("/merchants")
	merchants.Post("/", c.MerchantController.CreateMerchant)
	merchants.Get("/", c.MerchantController.GetAllMerchants)
	merchants.Get("/:id", c.MerchantController.GetMerchantById)
	merchants.Put("/:id", c.MerchantController.UpdateMerchant)
	merchants.Delete("/:id", c.MerchantController.DeleteMerchant)

	merchantProducts := api.Group("/merchant-products")
	merchantProducts.Post("/", c.MerchantProductController.CreateMerchantProduct)
	merchantProducts.Get("/:id", c.MerchantProductController.GetMerchantProductByID)
	merchantProducts.Get("/", c.MerchantProductController.GetMerchantProducts)
	merchantProducts.Get("/barcode/:barcode", c.MerchantProductController.GetMerchantProductByBarcode)
	merchantProducts.Put("/:id", c.MerchantProductController.UpdateMerchantProduct)
	merchantProducts.Delete("/:id", c.MerchantProductController.DeleteMerchantProduct)
	merchantProducts.Delete("/product/:product_id", c.MerchantProductController.DeleteAllProductMerchantProducts)
	merchantProducts.Get("/:product_id/total-stock", c.MerchantProductController.GetProductTotalStock)
}
