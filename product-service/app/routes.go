package app

import (
	_ "micro-warehouse/product-service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, container *Container) {
	api := app.Group("/api/v1")

	app.Get("/swagger/*", swagger.HandlerDefault)

	categories := api.Group("/categories")
	categories.Post("/", container.CategoryController.CreateCategory)
	categories.Get("/", container.CategoryController.GetAllCategories)
	categories.Get("/:id", container.CategoryController.GetCategoryByID)
	categories.Put("/:id", container.CategoryController.UpdateCategory)
	categories.Delete("/:id", container.CategoryController.DeleteCategory)

	products := api.Group("/products")
	products.Post("/", container.ProductController.CreateProduct)
	products.Get("/", container.ProductController.GetAllProducts)
	products.Get("/:id", container.ProductController.GetProductByID)
	products.Get("/barcode/:barcode", container.ProductController.GetProductByBarcode)
	products.Put("/:id", container.ProductController.UpdateProduct)
	products.Delete("/:id", container.ProductController.DeleteProduct)

	uploads := api.Group("/upload")
	uploads.Post("/product-image", container.UploadController.UploadProductImage)
	uploads.Post("/category-image", container.UploadController.UploadCategoryImage)
}
