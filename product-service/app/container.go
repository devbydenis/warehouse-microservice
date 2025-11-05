package app

import (
	"log"
	"micro-warehouse/product-service/configs"
	"micro-warehouse/product-service/controller"
	"micro-warehouse/product-service/database"
	"micro-warehouse/product-service/pkg/storage"
	"micro-warehouse/product-service/repository"
	"micro-warehouse/product-service/usecase"
)

type Container struct {
	ProductController  controller.ProductControllerInterface
	CategoryController controller.CategoryControllerInterface
	UploadController controller.UploadControllerInterface
}
func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectionPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	categoryRepo := repository.NewCategoryRepository(db.DB)
	categoryUseCase := usecase.NewCategoryUsecase(categoryRepo)
	categoryController := controller.NewCategoryController(categoryUseCase)

	productRepo := repository.NewProductRepository(db.DB)
	productUseCase := usecase.NewProductUsecase(productRepo)
	productController := controller.NewProductController(productUseCase)

	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewFileUploadHelper(supabaseStorage, *config)
	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		ProductController:  productController,
		CategoryController: categoryController,
		UploadController:  uploadController,
	}
}