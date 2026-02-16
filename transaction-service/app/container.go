package app

import (
	"log"
	"micro-warehouse/transaction-service/configs"
	"micro-warehouse/transaction-service/controller"
	"micro-warehouse/transaction-service/database"
	"micro-warehouse/transaction-service/pkg/httpclient"
	"micro-warehouse/transaction-service/pkg/midtrans"
	"micro-warehouse/transaction-service/pkg/rabbitmq"
	"micro-warehouse/transaction-service/repository"
	"micro-warehouse/transaction-service/usecase"
)

type Container struct {
	TransactionController controller.TransactionControllerInterface
}

func BuildContainer() *Container {
	cfg := configs.NewConfig()

	db, err := database.ConnectionPostgres(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// HTTP Client
	merchantClient := httpclient.NewMerchantClient(*cfg)
	productClient := httpclient.NewProductClient(*cfg)
	userClient := httpclient.NewUserCLient(*cfg)

	// RabbitMQ Client
	rabbitMQClient, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Midtrans Service
	midtransService := midtrans.NewMidtransService(cfg)

	transactionRepo := repository.NewTransactionRepository(db.DB)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, merchantClient, rabbitMQClient, productClient, userClient)
	transactionController := controller.NewTransactionController(transactionUsecase, midtransService)

	return &Container{
		TransactionController: transactionController,
	}
}