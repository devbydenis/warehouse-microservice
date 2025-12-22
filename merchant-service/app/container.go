package app

import (
	"micro-warehouse/merchant-service/configs"
	"micro-warehouse/merchant-service/controller"
	"micro-warehouse/merchant-service/database"
	"micro-warehouse/merchant-service/pkg/httpclient"
	"micro-warehouse/merchant-service/pkg/rabbitmq"
	"micro-warehouse/merchant-service/pkg/redis"
	"micro-warehouse/merchant-service/repository"
	"micro-warehouse/merchant-service/usecase"

	"github.com/gofiber/fiber/v2/log"
)

type Container struct {
	MerchantController controller.MerchantControllerInterface
	MerchantProductController controller.MerchantProductControllerInterface
}

func BuildContainer() *Container {
	cfg := configs.NewConfig()
	db, err := database.ConnectionPostgres(*cfg)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
	}

	redisClient := redis.NewRedisClient(*cfg)
	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL())
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %v", err)
	}

	userClient := httpclient.NewUserCLient(*cfg)
	cachedUserClient := httpclient.NewCachedUserClient(userClient, redisClient)
	warehouseClient := httpclient.NewWarehouseClient(*cfg)
	cachedWarehouseClient := httpclient.NewCachedWarehouseClient(warehouseClient, redisClient)
	productClient := httpclient.NewProductClient(*cfg)
	cachedProductClient := httpclient.NewCachedProductClient(productClient, redisClient)

	merchantRepo := repository.NewMerchantRepository(db.DB)
	merchantUsecase := usecase.NewMerchantUsecase(merchantRepo, cachedUserClient, cachedWarehouseClient, cachedProductClient)
	merchantController := controller.NewMerchantController(merchantUsecase)

	merchantProductRepo := repository.NewMerchantProductRepository(db.DB)
	merchantProductUsecase := usecase.NewMerchantProductUsecase(merchantProductRepo, cachedProductClient, cachedWarehouseClient, rabbitMQService)
	merchantProductController := controller.NewMerchantProductController(merchantProductUsecase)

	return &Container{
		MerchantController: merchantController,
		MerchantProductController: merchantProductController,
	}
}

