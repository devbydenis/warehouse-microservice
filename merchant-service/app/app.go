package app

import (
	"micro-warehouse/merchant-service/configs"
	"micro-warehouse/merchant-service/database"
	"micro-warehouse/merchant-service/pkg/rabbitmq"
	"micro-warehouse/merchant-service/repository"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	zerolog "github.com/rs/zerolog/log"
)

func RunServer() {
	cfg := configs.NewConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			zerolog.Printf("Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time} ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	container := BuildContainer()
	SetupRoutes(app, container)

	db, err := database.ConnectionPostgres(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	merchantProductRepo := repository.NewMerchantProductRepository(db.DB)
	stockCustomer, err := rabbitmq.NewStockConsumer(cfg.RabbitMQ.URL(), merchantProductRepo)
	if err != nil {
		log.Fatalf("Failed to create stock customer: %v", err)
	} else {
		go func() {
			ctx := context.Background()
			if err := stockCustomer.ConsumeStockReductionEvents(ctx); err != nil {
				log.Errorf("Failed to consume stock reduction events: %v", err)
			}
		}()
	}

	port := cfg.App.AppPort
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			log.Fatalf("Server port not specified")
		}
	}
	zerolog.Printf("Starting server on port: %s", port)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	zerolog.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Error during shutdown: %v", err)
	}

	zerolog.Printf("Server shutdown complete")

}
