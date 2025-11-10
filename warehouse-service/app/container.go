package app

import (
	"log"
	"micro-warehouse/warehouse-service/configs"
	"micro-warehouse/warehouse-service/controller"
	"micro-warehouse/warehouse-service/database"
	"micro-warehouse/warehouse-service/pkg/httpclient"
	"micro-warehouse/warehouse-service/repository"
	"micro-warehouse/warehouse-service/usecase"
)

type Container struct {
	WarehouseController        controller.WarehouseControllerInterface
	WarehouseProductController controller.WarehouseProductControllerInterface
}

func BuildContainer() *Container {
	config := configs.NewConfig()
	db, err := database.ConnectionPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	warehouseRepo := repository.NewWarehouseRepository(db.DB)
	warehouseUsecase := usecase.NewWarehouseUsecase(warehouseRepo)
	warehouseController := controller.NewWarehouseController(warehouseUsecase)

	warehouseProductRepo := repository.NewWarehouseProductRepository(db.DB)
	productClient := httpclient.NewProductClient(*config)
	warehouseProductUsecase := usecase.NewWarehouseProductUsecase(warehouseProductRepo, productClient)
	warehouseProductController := controller.NewWarehouseProductController(warehouseProductUsecase)
	
	return &Container{
		WarehouseController:        warehouseController,
		WarehouseProductController: warehouseProductController,
	}
}
