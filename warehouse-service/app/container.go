package app

import (
	"log"
	"time"
	
	"micro-warehouse/warehouse-service/configs"
	"micro-warehouse/warehouse-service/controller"
	"micro-warehouse/warehouse-service/database"
	"micro-warehouse/warehouse-service/pkg/httpclient"
	"micro-warehouse/warehouse-service/pkg/redis"
	"micro-warehouse/warehouse-service/pkg/storage"
	"micro-warehouse/warehouse-service/repository"
	"micro-warehouse/warehouse-service/usecase"
)

type Container struct {
	WarehouseController        controller.WarehouseControllerInterface
	WarehouseProductController controller.WarehouseProductControllerInterface
	UploadController           controller.UploadControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectionPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	productClient := httpclient.NewProductClient(*config)
	redisClient := redis.NewRedisClient(*config)
	cachedProdctClient := httpclient.NewCachedProductClient(productClient, redisClient, 1*time.Hour)

	warehouseRepo := repository.NewWarehouseRepository(db.DB)
	warehouseUsecase := usecase.NewWarehouseUsecase(warehouseRepo)
	warehouseController := controller.NewWarehouseController(warehouseUsecase)

	warehouseProductRepo := repository.NewWarehouseProductRepository(db.DB)
	warehouseProductUsecase := usecase.NewWarehouseProductUsecase(warehouseProductRepo, cachedProdctClient)
	warehouseProductController := controller.NewWarehouseProductController(warehouseProductUsecase)

	supabaseStorage := storage.NewSupabaseStorage(*config)
	fileUploadHelper := storage.NewFileUploadHelper(supabaseStorage, *config)
	uploadController := controller.NewUploadController(fileUploadHelper)

	return &Container{
		WarehouseController:        warehouseController,
		WarehouseProductController: warehouseProductController,
		UploadController:           uploadController,
	}
}
